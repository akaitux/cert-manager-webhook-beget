{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "beget-webhook.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "beget-webhook.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "beget-webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "beget-webhook.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "beget-webhook.fullname" .) }}
{{- end -}}

{{- define "beget-webhook.rootCAIssuer" -}}
{{ printf "%s-ca" (include "beget-webhook.fullname" .) }}
{{- end -}}

{{- define "beget-webhook.rootCACertificate" -}}
{{ printf "%s-ca" (include "beget-webhook.fullname" .) }}
{{- end -}}

{{- define "beget-webhook.servingCertificate" -}}
{{ printf "%s-webhook-tls" (include "beget-webhook.fullname" .) }}
{{- end -}}
