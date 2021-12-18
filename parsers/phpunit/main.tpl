<testsuites>
  {{- range $l1suite := .testsuite -}}
    {{- range $l2suite := .testsuite -}}
      {{- range $l3suite := .testsuite -}}
        <testsuite {{ attributes . }}>
          {{ range $test := .testcase }}
            {{ template "toXML" $test }}
          {{ end }}
        </testsuite>
      {{ end }}
    {{ end }}
  {{ end }}
</testsuites>