<img align="right" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">

# mapreduce

[English](readme.md) | [简体中文](readme-cn.md) | 한국어

## MapReduce가 필요한 이유

실제 비즈니스 시나리오에서는 서로 다른 RPC 서비스에서 속성을 가져와 복잡한 객체를 조립해야 하는 경우가 많습니다.

예를 들어 상품 상세 정보를 조회한다고 해봅시다.

1. 상품 서비스 - 상품 속성 조회
2. 재고 서비스 - 재고 속성 조회
3. 가격 서비스 - 가격 속성 조회
4. 마케팅 서비스 - 마케팅 속성 조회

직렬 호출이라면 RPC 호출 횟수에 따라 응답 시간이 선형적으로 증가하므로, 일반적으로 응답 시간을 최적화하기 위해 직렬 호출을 병렬 호출로 바꿉니다.

단순한 시나리오에서는 `WaitGroup`만으로도 요구 사항을 충족할 수 있습니다. 하지만 RPC 호출이 반환한 데이터를 검증하거나, 데이터를 처리하거나, 데이터를 집계해야 한다면 어떻게 해야 할까요? Go 표준 라이브러리에는 이런 도구가 없습니다(Java에는 CompletableFuture가 제공됩니다). 그래서 우리는 MapReduce 아키텍처를 기반으로 프로세스 내부 데이터 배치 처리를 위한 MapReduce 동시성 도구를 구현했습니다.

## 설계 아이디어

동시성 도구가 필요한 비즈니스 시나리오를 작성자의 관점에서 정리해봅시다.

1. 상품 상세 조회: 여러 서비스를 동시에 호출해 상품 속성을 조합하고, 호출 오류가 발생하면 즉시 종료할 수 있어야 합니다.
2. 상품 상세 페이지에서 사용자 쿠폰 자동 추천: 쿠폰을 동시에 검증하고, 검증에 실패한 쿠폰은 자동으로 제외하며, 나머지를 모두 반환할 수 있어야 합니다.

위 시나리오는 모두 입력 데이터를 처리한 뒤 정제된 데이터를 출력하는 과정입니다. 데이터 처리에는 아주 고전적인 비동기 패턴인 생산자-소비자 패턴이 있습니다. 따라서 데이터 배치 처리의 생명 주기를 추상화하면 대략 세 단계로 나눌 수 있습니다.

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/mapreduce-serial-en.png" width="500">

1. 데이터 생성(generator)
2. 데이터 처리(mapper)
3. 데이터 집계(reducer)

데이터 생성은 필수 단계이고, 데이터 처리와 데이터 집계는 선택 단계입니다. 데이터 생성과 처리는 동시 호출을 지원하며, 데이터 집계는 기본적으로 순수 메모리 작업이므로 단일 고루틴으로 처리할 수 있습니다.

서로 다른 데이터 처리 단계가 서로 다른 고루틴에서 수행되므로, 고루틴 간 통신을 위해 채널을 사용하는 것이 자연스럽습니다.

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/mapreduce-en.png" width="500">

언제든 프로세스를 종료하려면 어떻게 해야 할까요?

간단합니다. 고루틴에서 채널 또는 전달된 context의 완료 신호를 감시하면 됩니다.

## 간단한 예시

동시성을 시뮬레이션하며 제곱합을 계산합니다.

```go
package main

import (
    "fmt"
    "log"

    "github.com/zeromicro/go-zero/core/mr"
)

func main() {
    val, err := mr.MapReduce(func(source chan<- int) {
        // generator
        for i := 0; i < 10; i++ {
            source <- i
        }
    }, func(i int, writer mr.Writer[int], cancel func(error)) {
        // mapper
        writer.Write(i * i)
    }, func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
        // reducer
        var sum int
        for i := range pipe {
            sum += i
        }
        writer.Write(sum)
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("result:", val)
}
```

더 많은 예제: [https://github.com/zeromicro/zero-examples/tree/main/mapreduce](https://github.com/zeromicro/zero-examples/tree/main/mapreduce)

## 별을 눌러주세요! ⭐

이 프로젝트가 마음에 들거나 학습 또는 자체 솔루션을 시작하는 데 사용 중이라면 star를 눌러주세요. 감사합니다!
