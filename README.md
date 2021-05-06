# Test Results

## CLI Usage

### Job Level

#### Publishing a single report from a Job.

``` bash
$ test-results publish junit.xml

* Detected ExUnit JUnit test results
* Processing report, saving results to /tmp/test-results103206663
* Pushing raw report to artifacts:
  artifact push job junit.xml -d test-results/raw/junit.xml
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/report.json
```

#### Manually selecting the parser for test results.

``` bash
$ test-results publish junit.xml --parser ExUnit

* Using ExUnit parser to process results
* Processing report, saving results to /tmp/test-results103206663
* Pushing raw report to artifacts:
  artifact push job junit.xml -d test-results/raw/junit.xml
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/report.json
```

#### Asking the test-results CLI to not publish raw reports.

``` bash
$ test-results publish junit.xml --parser ExUnit --no-raw

* Detected ExUnit JUnit test results
* Processing report, saving results to /tmp/test-results103206663
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/report.json
```

#### Asking the test-results CLI to set an expiration date for published artifacts.

``` bash
$ test-results publish junit.xml --parser ExUnit --expire-in 7d

* Detected ExUnit JUnit test results
* Processing report, saving results to /tmp/test-results103206663
* Pushing raw report to artifacts:
  artifact push job junit.xml -d test-results/raw/junit.xml --expire-in 7d
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/report.json --expire-in 7d
```

#### Publishing multiple test-results from a folder:

``` bash
$ test-results publish results

* 3 JUnit XML reports in results/
  - Using ExUnit parser for results/a.xml 
  - Using Mocha parser for results/b.xml 
* Processing reports, saving results to /tmp/test-results103206663
* Pushing raw report to artifacts:
  artifact push job results/ -d test-results/raw/
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/report.json
```

#### Publishing multiple test-results with manually selecting a parser:

WIP.

``` bash
$ test-results publish results/a.xml --parser ExUnit

* Detected ExUnit JUnit test results
* Processing reports, saving results to /tmp/test-results103206663
* Pushing raw report to artifacts:
  artifact push job results/a.xml -d test-results/junit.xml
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/junit.json

$ test-results publish results/b.xml --parser Mocha

* Found existing test result report in artifacts. Pulling it down for extension.
  artifact pull job results/junit.json -d test-results/junit.json
* Detected ExUnit JUnit test results
* Processing reports, saving results to /tmp/test-results103206663
* Pushing raw report to artifacts:
  artifact push job results/a.xml -d test-results/junit.xml
* Pushing processed report to artifacts:
  artifact push job /tmp/test-results103206663 -d test-results/junit.json
```

### Pipeline level

WIP

``` bash
$ test-results gen-pipeline-report
* Fetching job level reports
  - artifact pull job test-results/report.json -d /tmp/test-results/67bb1901-4823-4e01-8af5-d992dc3b6792.json
  - artifact pull job test-results/report.json -d /tmp/test-results/e8ece834-c9e5-4e47-867c-26af7b0c82f8.json
  - artifact pull job test-results/report.json -d /tmp/test-results/bc8c9bd2-16d4-4321-b376-719a0531a07b.json
* Merging reports, saving result to report.json
* Pushing report to artifacts:
  artifact push workflow report.json -d test-results/$SEMAPHORE_PIPELINE_ID.json
```

How to pull job level reports with Artifacts?

Where to store reports:

```
Approach #1:
------------

Job A:

  $ test-results publish junit.xml
  artifact push workflow test-results/$PIPELINE_ID/$JOB_ID.json

Job B:

  $ test-results publish junit.xml
  artifact push workflow test-results/$PIPELINE_ID/$JOB_ID.json
  
After Job:
  
  $ test-results gen-pipeline-report
  artifacts pull workflow test-results/$PIPELINE_ID --destination /tmp/test-results
  Merging....
  artifacts push workflow test-results/$PIPELINE_ID.json
  
Approach #2:
-----------


Job A:

  $ test-results publish junit.xml
  artifact push project test-results/$WORKFLOW_ID/$PIPELINE_ID/$JOB_ID.json

Job B:

  $ test-results publish junit.xml
  artifact push project test-results/$WORKFLOW_ID/$PIPELINE_ID/$JOB_ID.json
  
After Job:
  
  $ test-results gen-pipeline-report
  artifacts pull project test-results/$WORKFLOW_ID/$PIPELINE_ID --destination /tmp/test-results
  Merging....
  artifacts push project test-results/$WORKFLOW_ID/$PIPELINE_ID.json

Approach #3:
------------

Job A:

  $ test-results publish junit.xml
  artifact push job test-results/report.json

Job B:

  $ test-results publish junit.xml
  artifact push job test-results/report.json
  
After job:
  # generate this file from Zebra:
  
  $ cat ~/jobs.txt
  67bb1901-4823-4e01-8af5-d992dc3b6792
  e8ece834-c9e5-4e47-867c-26af7b0c82f8
  746203fa-b635-46ca-b4fa-d4814f274069
  
  $ test-results gen-pipeline-report
```
