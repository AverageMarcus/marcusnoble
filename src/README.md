<div align="center">

# Hi 👋, I'm Marcus

<sup>**Pronounced:** /ˈmɑːkəs/ ⚡️ **Pronouns:** He/Him/His</sup>

{{ .intro | html }}

## 🌐 Find me around the web

| {{ range .social }}<a href="{{ .url }}" rel="me" title="{{ .title }}">{{ .name | html }}</a> | {{ end }}

{{ if gt (len .events) 0 }}
## 🗓 Upcoming Events

{{ range .events }}
<div>{{ .humanDate }}</div>
<div>

[**{{.eventName}}**]({{ .url }})

</div>
{{ range .details }}
<strong>

{{ .name | html }}{{ if .type }} - {{ .type }}{{ end}}

</strong>
{{ end }}
✨✨✨
{{- end }}

</div>

{{ end }}

{{ if gt (len .projects) 0 }}
## 💻 My Open Source Projects

All my Open Source projects can be found on my <a href="https://github.com/AverageMarcus">GitHub</a> profile (as well as my personal <a href="https://git.cluster.fun">Gitea</a> instance, <a href="https://gitlab.com/AverageMarcus">GitLab</a>, <a href="https://codeberg.org/AverageMarcus">Codeberg</a> and <a href="https://bitbucket.org/AverageMarcus/workspace/projects/PROJ">BitBucket</a>).

Below are a selection of highlights.

{{ range .projects }}
[**{{ .name }}**]({{ .url }}) - {{ .description }} [{{ join .languages "name" ", " }}]
{{ end }}
{{ end }}

{{ if gt (len .badges) 0 }}
## 🎖️ Badges

{{ range .badges }}
<a href="{{ .url }}" target="_blank">
  <img src="{{ .img }}" width="110px" />
</a>
{{ end }}
{{ end }}
