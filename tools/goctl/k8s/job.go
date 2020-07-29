package k8s

// 无环境区分
var jobTmeplate = `apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: {{.name}}
  namespace: {{.namespace}}
spec:
  successfulJobsHistoryLimit: {{.successfulJobsHistoryLimit}}
  schedule: "{{.schedule}}"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: {{.name}}
            image: registry-vpc.cn-hangzhou.aliyuncs.com/{{.namespace}}/
            env:
            - name: aliyun_logs_k8slog
              value: "stdout"
            - name: aliyun_logs_k8slog_tags
              value: "stage={{.env}}"
            - name: aliyun_logs_k8slog_format
              value: "json"
            resources:
              limits:
                cpu: {{.limitCpu}}m
                memory: {{.limitMem}}Mi
              requests:
                cpu: {{.requestCpu}}m
                memory: {{.requestMem}}Mi
            command:
            - ./{{.serviceName}}
            - -f
            - ./{{.name}}.json
            volumeMounts:
            - name: timezone
              mountPath: /etc/localtime
          imagePullSecrets:
          - name: {{.namespace}}
          restartPolicy: OnFailure
          volumes:
          - name: timezone
            hostPath:
              path: /usr/share/zoneinfo/Asia/Shanghai`
