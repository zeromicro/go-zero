package sarama

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

type FinishReason string

const (
	Unexpected                   FinishReason = "Unexpected"
	OnConsumeLoopInitReturnFalse FinishReason = "OnConsumeLoopInitReturnFalse"
	MessagesChannelClosed        FinishReason = "MessagesChannelClosed"
	OnConsumeMessageReturnFalse  FinishReason = "OnConsumeMessageReturnFalse"
	OnConsumeErrorReturnFalse    FinishReason = "OnConsumeErrorReturnFalse"
)

type (
	MessageHandler interface {
		// OnKNameConfirmed 根据配置确认出 从哪个kafka集群读消息
		OnKNameConfirmed(k string)

		// OnConsumeLoopInit 返回false直接结束并close, 不开始消费
		OnConsumeLoopInit(currentBeginInclusive, offset, currentEndExclusive int64) (continueConsuming bool)

		// OnConsumeMessage 返回false中止消费, consumerLoop不会再调用OnRecvMessage()方法 并开始执行优雅关闭逻辑
		OnConsumeMessage(ctx context.Context, msg *types.Message, highWaterMarkOffset int64) (continueConsuming bool)

		// OnConsumeError 返回true表示继续消费, 返回false优雅关闭consumer
		OnConsumeError(consumerErr *ConsumerError, highWaterMarkOffset int64) (continueConsuming bool)

		// OnConsumeLoopFinish consumer协程结束时回调该函数, 此时partitionConsumer已经close完成, map已经清理, wait group也已经-1. 可以新开另一个consumer
		OnConsumeLoopFinish(finishReason FinishReason)
	}

	// UnsafeConsumer is a wrapper of sarama.Consumer.
	UnsafeConsumer interface {
		Consumer
		SpawnPartitionConsumer(topic string, partition int32, offset, begin, end int64, cHandler MessageHandler) (UnsafePartitionConsumer, error)
		ClosePartitionConsumer(ctx context.Context, topic string, partition int32) bool
		SyncClosePartitionConsumer(ctx context.Context, topic string, partition int32) bool
		CloseAllPartitionConsumers(ctx context.Context)
	}

	// UnsafePartitionConsumer is a wrapper of sarama.PartitionConsumer.
	UnsafePartitionConsumer interface {
		PartitionConsumer
		AsyncClose()
	}

	unsafeConsumer struct {
		*consumer
		k string
	}

	unsafePartitionConsumer struct {
		unsafeConsumer     *unsafeConsumer
		spc                sarama.PartitionConsumer
		cHandler           MessageHandler
		done               chan struct{}
		once               sync.Once
		offset, begin, end int64
		uniq               string
	}
)

func newUnsafeConsumerFromClient(c *Client) (*unsafeConsumer, error) {
	con, err := newConsumerFromClient(c)
	con.sharedClient = c
	sc := &unsafeConsumer{
		consumer: con,
	}
	return sc, err
}

// SpawnPartitionConsumer warps sarama.consumer.ConsumePartition with handlers,
func (c *unsafeConsumer) SpawnPartitionConsumer(topic string, partition int32, offset, begin, end int64, handler MessageHandler) (UnsafePartitionConsumer, error) {
	p := &unsafePartitionConsumer{
		unsafeConsumer: c,
		cHandler:       handler,
		offset:         offset,
		begin:          begin,
		end:            end,
		done:           make(chan struct{}),
		uniq:           fmt.Sprintf("%s-%d", topic, partition),
	}
	_, exists := c.sharedClient.spawnedPartitionConsumers.Load(p.uniq)
	if exists {
		// 同一个client不能多次消费同一队列!!!
		return nil, fmt.Errorf("consumer_%s already eixsts", p.uniq)
	}
	pc, err := c.saramaConsumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, err
	}
	p.spc = pc
	//注释锁之后，防止并发调用SpawnPartitionConsumer。在store之前再次检查是否loaded
	_, exists = c.sharedClient.spawnedPartitionConsumers.LoadOrStore(p.uniq, p)
	if exists {
		_ = p.Close()
		// 同一个client不能多次消费同一队列!!!
		return nil, fmt.Errorf("consumer_%s already eixsts", p.uniq)
	}
	p.setupCallbacks()
	logx.Infof("consumer partition started, name: %s, topic: %s, partition: %d, startOffset: %d", p.unsafeConsumer.name, topic, partition, offset)
	return p, nil
}

func (c *unsafeConsumer) Close() error {
	if c.saramaClient == nil {
		return nil
	}

	var err error
	c.once.Do(func() {
		c.CloseAllPartitionConsumers(context.TODO())
		if c.closeClient {
			err = c.saramaClient.Close()
		}
	})

	return err
}

func (c *unsafeConsumer) CloseAllPartitionConsumers(ctx context.Context) {
	logx.Infof("CloseAllPartitionConsumers of %s: trigger %s async close", c.k, c.name)
	c.sharedClient.spawnedPartitionConsumers.Range(func(uniq, pc any) bool {
		c.closePartitionConsumer(ctx, uniq.(string), false)
		return true // 一直循环
	})
}

func (c *unsafeConsumer) ClosePartitionConsumer(ctx context.Context, topic string, partition int32) bool {
	uniq := fmt.Sprintf("%s-%d", topic, partition)
	return c.closePartitionConsumer(ctx, uniq, false)
}

func (c *unsafeConsumer) SyncClosePartitionConsumer(ctx context.Context, topic string, partition int32) bool {
	uniq := fmt.Sprintf("%s-%d", topic, partition)
	return c.closePartitionConsumer(ctx, uniq, true)
}

func (c *unsafeConsumer) closePartitionConsumer(ctx context.Context, uniq string, syncClose bool) bool {
	pc, exists := c.sharedClient.spawnedPartitionConsumers.Load(uniq)
	if exists {
		c.sharedClient.spawnedPartitionConsumers.Delete(uniq)
		if syncClose {
			logx.Infof("ClosePartitionConsumer_%s: trigger %s sync close", c.k, uniq)
			_ = pc.(UnsafePartitionConsumer).Close()
		} else {
			logx.Infof("ClosePartitionConsumer_%s: trigger %s async close", c.k, uniq)
			pc.(UnsafePartitionConsumer).AsyncClose()
		}
		return true
	} else {
		return false
	}
}

func (p *unsafePartitionConsumer) setupCallbacks() {
	p.unsafeConsumer.sharedClient.cwg.Add(1) // 新创建一个consumer的消费协程, 引用计数+1
	go func() {
		finishReason := Unexpected
		defer func() {
			p.unsafeConsumer.sharedClient.spawnedPartitionConsumers.Delete(p.uniq)
			err := p.Close()
			if err != nil {
				logx.Errorw(fmt.Sprintf("baseConsumer_%s.Close err %v", p.uniq, err),
					logx.Field("brokers", p.unsafeConsumer.name),
				)
			}
			p.unsafeConsumer.sharedClient.cwg.Done() // ConsumerCallbackLoop协程结束后, 引用计数-1
			p.cHandler.OnConsumeLoopFinish(finishReason)
		}()

		if p.cHandler.OnConsumeLoopInit(p.begin, p.offset, p.end) == false {
			err := p.Close()
			if err != nil {
				logx.Errorw(fmt.Sprintf("[init()->false] err on partitionConsumer_%s.Close() %v", p.uniq, err),
					logx.Field("brokers", p.unsafeConsumer.name),
				)
			}

			finishReason = OnConsumeLoopInitReturnFalse
			return // 直接正常结束
		}

		successChClosed := false
	PartitionConsumerProcessorLoop:
		for {
			select {
			case msg, ok := <-p.spc.Messages():
				if ok {
					if !p.handleMessage(msg, p.HighWaterMarkOffset()) {
						err := p.Close()
						if err != nil {
							logx.Errorw(fmt.Sprintf("[msg()->false] err on partitionConsumer_%s.Close(),topic:%s,partition:%d,offset:%d err:%v", p.uniq, msg.Topic, msg.Partition, msg.Offset, err),
								logx.Field("brokers", p.unsafeConsumer.name),
							)
						}
						finishReason = OnConsumeMessageReturnFalse
						return // 提前正常结束
					}

				} else if !successChClosed {
					logx.Infof("partitionConsumer_%s[%s] Messages channel closed", p.uniq, p.uniq)
					successChClosed = true // 标记为已经处理过close
				}

			case consumerErr, ok := <-p.spc.Errors():
				if ok {
					//后close errors 所以 建议ok判断放在errors上?
					logx.Errorw(fmt.Sprintf("err on partitionConsumer_%s: tp:%s, err:%v", p.uniq, p.uniq, consumerErr),
						logx.Field("brokers", p.unsafeConsumer.name),
					)
					// sarama内部的实现是会自己重试consumer error, 所以这里无需panic

					// 注册callback后, 由上层业务自己判断是否需要继续consuming
					if !p.cHandler.OnConsumeError(consumerErr, p.HighWaterMarkOffset()) {
						err := p.Close()
						if err != nil {
							logx.Errorw(fmt.Sprintf("[err()->false] err on partitionConsumer_%s.Close(),tp:%s, err:%v", p.uniq, p.uniq, consumerErr),
								logx.Field("brokers", p.unsafeConsumer.name),
							)
						}
						finishReason = OnConsumeErrorReturnFalse
						return // 提前正常结束
					} else {
						logx.Debugf("partitionConsumer_%s continue consuming tp:%s", p.uniq, p.uniq)
					}
				} else {
					logx.Infof("partitionConsumer_%s Errors channel closed  tp:%s", p.uniq, p.uniq)
					finishReason = MessagesChannelClosed
					// sarama的consumer底层是先close(successes) 再close(errors)
					break PartitionConsumerProcessorLoop // 需要后置判断是否panic
				}
			}
		}

		// 主动调用AsyncClose的话, 会先将entry移除掉
		_, partitionConsumerStillExists := p.unsafeConsumer.sharedClient.spawnedPartitionConsumers.Load(p.uniq)

		// consumer loop结束后, 判断下是用户发起的close 还是底层error触发的close 决定是否直接panic
		if partitionConsumerStillExists && !p.unsafeConsumer.sharedClient.isExiting {
			panic(fmt.Errorf("unexpected partitionConsumer_%s[%s] consumer loop end", p.uniq, p.uniq))
		}
	}()
}

func (p *unsafePartitionConsumer) handleMessage(message *sarama.ConsumerMessage, highWaterMarkOffset int64) bool {
	return innerHandleMessage(p.unsafeConsumer.consumer, message, true, func(ctx context.Context, msg *types.Message) bool {
		return p.cHandler.OnConsumeMessage(ctx, msg, highWaterMarkOffset)
	})
}

func (p *unsafePartitionConsumer) Close() error {
	if p.spc == nil {
		return nil
	}

	p.once.Do(func() {
		logx.Infof("partitionConsumer sync close, tp:%s", p.uniq)
		close(p.done)
		err := p.spc.Close()
		if err != nil {
			logx.Errorw(fmt.Sprintf("[init()->false] err on partitionConsumer_%s.Close() %v", p.uniq, err.Error()),
				logx.Field("brokers", p.unsafeConsumer.name),
			)
		}
	})

	return nil
}

func (p *unsafePartitionConsumer) AsyncClose() {
	if p.spc == nil {
		return
	}

	p.once.Do(func() {
		logx.Infof("partitionConsumer async close tp: %s", p.uniq)
		close(p.done)
		p.spc.AsyncClose()
	})
}

func (p *unsafePartitionConsumer) HighWaterMarkOffset() int64 {
	return p.spc.HighWaterMarkOffset()
}
