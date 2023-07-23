apiVersion: batch/v1
kind: CronJob
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  successfulJobsHistoryLimit: {{.SuccessfulJobsHistoryLimit}}
  schedule: "{{.Schedule}}"
  jobTemplate:
    spec:
      template:
        spec:{{if .ServiceAccount}}
          serviceAccountName: {{.ServiceAccount}}{{end}}
	      {{end}}containers:
          - name: {{.Name}}
            image: # todo image url
            resources:
              requests:
                cpu: {{.RequestCpu}}m
                memory: {{.RequestMem}}Mi
              limits:
                cpu: {{.LimitCpu}}m
                memory: {{.LimitMem}}Mi
            command:
            - ./{{.ServiceName}}
            - -f
            - ./{{.Name}}.yaml
            volumeMounts:
            - name: timezone
              mountPath: /etc/localtime
          imagePullSecrets:
          - name: # registry secret, if no, remove this
          restartPolicy: OnFailure
          volumes:
          - name: timezone
            hostPath:
              path: /usr/share/zoneinfo/Asia/Shanghai
