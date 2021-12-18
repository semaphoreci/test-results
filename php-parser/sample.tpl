<testsuites name="Customized example">
  {{ range $l1suite := .testsuite }}
    {{ range $l2suite := .testsuite }}
      {{ range $l3suite := .testsuite }}
        <testsuite 
          name='{{ "name" | field $l2suite }} ➡️  {{ "name" | field $l3suite }}'
          tests='{{ "tests" | field $l3suite }}'
          assertions='{{ "assertions" | field $l3suite }}'
          errors='{{ "errors" | field $l3suite }}'
          warnings='{{ "warnings" | field $l3suite }}'
          failures='{{ "failures" | field $l3suite }}'
          skipped='{{ "skipped" | field $l3suite }}'
          time='{{ "time" | field $l3suite }}'
        >
          {{ range $test := .testcase }}
            {{ template "toXML" $test }}
          {{ end }}
        </testsuite>
      {{ end }}
    {{ end }}
  {{ end }}
</testsuites>