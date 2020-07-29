package k8s

var rmqSyncTmeplate = `apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: {{.name}}
  namespace: {{.namespace}}
  labels:
    app: {{.name}}
spec:
  replicas: {{.replicas}}
  revisionHistoryLimit: {{.revisionHistoryLimit}}
  selector:
    matchLabels:
      app: {{.name}}
  template:
    metadata:
      labels:
        app: {{.name}}
    spec:{{if .envIsDev}}
      terminationGracePeriodSeconds: 60{{end}}
      containers:
      - name: {{.name}}
        image: registry-vpc.cn-hangzhou.aliyuncs.com/{{.namespace}}/
        lifecycle:
          preStop:
            exec:
              command: ["sh","-c","sleep 5"]
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
      volumes:
        - name: timezone
          hostPath:
            path: /usr/share/zoneinfo/Asia/Shanghai{{if .envIsPreOrPro}}

---
apiVersion: v1
kind: Service
metadata:
  name: {{.name}}-svc
  namespace: {{.namespace}}
spec:
  selector:
    app: {{.name}}
  sessionAffinity: None
  type: ClusterIP
  clusterIP: None{{end}}`
