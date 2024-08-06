{{/*
Expand the name of the chart.
*/}}
{{- define "notification-manager.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "notification-manager.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "notification-manager.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "notification-manager.labels" -}}
helm.sh/chart: {{ include "notification-manager.chart" . }}
{{ include "notification-manager.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "notification-manager.selectorLabels" -}}
app.kubernetes.io/name: {{ include "notification-manager.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
#########################
# notification-manager  grpc#
#########################
*/}}
{{/* Labels */}}
{{- define "notification-manager-grpc.labels" -}}
helm.sh/chart: {{ include "notification-manager.chart" . }}
{{ include "notification-manager.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "notification-manager-grpc.selectorLabels" -}}
app.kubernetes.io/name: {{ include "notification-manager.name" . }}-grpc
app.kubernetes.io/instance: {{ .Release.Name }}-grpc
{{- end }}

{{/*
#########################
# notification-manager worker master #
#########################
*/}}

{{/* Labels */}}
{{- define "notification-manager-worker-master.labels" -}}
helm.sh/chart: {{ include "notification-manager.chart" . }}
{{ include "notification-manager-worker-master.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "notification-manager-worker-master.selectorLabels" -}}
app.kubernetes.io/name: {{ include "notification-manager.name" . }}-worker-master
app.kubernetes.io/instance: {{ .Release.Name }}-worker-master
{{- end }}

{{/*
#########################
# notification-manager worker email #
#########################
*/}}

{{/* Labels */}}
{{- define "notification-manager-worker-email.labels" -}}
helm.sh/chart: {{ include "notification-manager.chart" . }}
{{ include "notification-manager-worker-email.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "notification-manager-worker-email.selectorLabels" -}}
app.kubernetes.io/name: {{ include "notification-manager.name" . }}-worker-email
app.kubernetes.io/instance: {{ .Release.Name }}-worker-email
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "notification-manager.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "notification-manager.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
