package k8s

var apiRpcTmeplate = `apiVersion: apps/v1beta2
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
        ports:
        - containerPort: {{.port}}
        readinessProbe:
          tcpSocket:
            port: {{.port}}
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          tcpSocket:
            port: {{.port}}
          initialDelaySeconds: 15
          periodSeconds: 20
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
            path: /usr/share/zoneinfo/Asia/Shanghai

---
apiVersion: v1
kind: Service
metadata:
  name: {{.name}}-svc
  namespace: {{.namespace}}
spec:
  ports:
    - nodePort: 3{{.port}}
      port: {{.port}}
      protocol: TCP
      targetPort: {{.port}}
  selector:
    app: {{.name}}
  sessionAffinity: None
  type: NodePort{{if .envIsPreOrPro}}

---
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: {{.name}}-hpa-c
  namespace: {{.namespace}}
  labels:
    app: {{.name}}-hpa-c
spec:
  scaleTargetRef:
    apiVersion: apps/v1beta1
    kind: Deployment
    name: di-api
  minReplicas: {{.minReplicas}}
  maxReplicas: {{.maxReplicas}}
  metrics:
  - type: Resource
    resource:
      name: cpu
      targetAverageUtilization: 80

---
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: {{.name}}-hpa-m
  namespace: {{.namespace}}
  labels:
    app: {{.name}}-hpa-m
spec:
  scaleTargetRef:
    apiVersion: apps/v1beta1
    kind: Deployment
    name: {{.name}}
  minReplicas: {{.minReplicas}}
  maxReplicas: {{.maxReplicas}}
  metrics:
  - type: Resource
    resource:
      name: memory
      targetAverageUtilization: 80{{end}}`
