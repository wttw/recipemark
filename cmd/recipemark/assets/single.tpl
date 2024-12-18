<html>
<head>
    <title>{{ .Name }}</title>
    <meta name="title" on
</head>
<body>
<div class="preamble">
    {{ range .Chunks -}}
        {{ if eq .Section "romance" }}{{ .Content }}{{ end -}}
    {{- end }}
</div>
<div itemscope itemtype="https://schema.org/Recipe">
    <div class="ingredients">
        {{ range .Chunks -}}
            {{ if eq .Section "ingredients" }}{{ .Content }}{{ end -}}
        {{- end }}
    </div>
    <div class="method">
        {{ range .Chunks -}}
            {{ if eq .Section "method" }}{{ .Content }}{{ end -}}
        {{- end }}
    </div>
</div>
</body>
</html>