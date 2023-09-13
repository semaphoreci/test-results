package parsers

import (
	"bytes"
	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Embedded_ParseTestSuites(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
	    <testsuites>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testsuite name="zap" id="4321">
					<testcase name="baz">
					</testcase>
				</testsuite>
				<testsuite name="zup" id="54321">
					<testcase name="bar">
					</testcase>
				</testsuite>
				<testsuite name="testNumSuccess(boolean)" tests="2" failures="0" errors="1" disabled="0"
					skipped="0" package="">
					<properties />
					<testcase name="[1] testVirtualMetrics=true"
						classname="io.testcompany.ZedCounterAdminTest" time="0">
						<error message="expected: &lt;true&gt; but was: &lt;false&gt;"
							type="org.opentest4j.AssertionFailedError"><![CDATA[org.opentest4j.AssertionFailedError: expected: <true> but was: <false>
			at org.junit.jupiter.api.AssertionFailureBuilder.build(AssertionFailureBuilder.java:151)
			at org.junit.jupiter.api.AssertionFailureBuilder.buildAndThrow(AssertionFailureBuilder.java:132)
			at org.junit.jupiter.api.AssertTrue.failNotTrue(AssertTrue.java:63)
			at org.junit.jupiter.api.AssertTrue.assertTrue(AssertTrue.java:36)
			at org.junit.jupiter.api.AssertTrue.assertTrue(AssertTrue.java:31)
			at org.junit.jupiter.api.Assertions.assertTrue(Assertions.java:180)
			at io.testcompany.ZedCounterAdminTest.testNumSuccess(ZedCounterAdminTest.java:35)
			at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
			at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:77)
			at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
			at java.base/java.lang.reflect.Method.invoke(Method.java:568)
			at org.junit.platform.engine.support.hierarchical.NodeTestTask.lambda$executeRecursively$9(NodeTestTask.java:139)
			at org.junit.platform.engine.support.hierarchical.ThrowableCollector.execute(ThrowableCollector.java:73)
			at org.junit.platform.engine.support.hierarchical.NodeTestTask.executeRecursively(NodeTestTask.java:138)
			at org.junit.platform.engine.support.hierarchical.NodeTestTask.execute(NodeTestTask.java:95)
			at java.base/java.util.ArrayList.forEach(ArrayList.java:1511)
		]]></error>
					</testcase>
				<testcase name="[2] testVirtualMetrics=false"
					classname="io.testcompany.ZedCounterAdminTest" time="0.1"></testcase>
			</testsuite>
		</testsuite>
	</testsuites>
	`))

	path := fileloader.Ensure(reader)

	e := NewEmbedded()
	testResults := e.Parse(path)
	assert.Equal(t, "Suite", testResults.Name)
	assert.Equal(t, "embedded", testResults.Framework)
	assert.Equal(t, "c5bec5ae-e57f-3dac-98fa-825a5a2cfd55", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	require.Equal(t, 4, len(testResults.Suites))

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "60e78e69-056c-3806-b86f-19fc2f9b5124", testResults.Suites[0].ID)
	require.Len(t, testResults.Suites[0].Tests, 1)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)

	assert.Equal(t, "foo\\zap", testResults.Suites[1].Name)
	assert.Equal(t, "427236d0-f6f4-3483-a023-95c7fce844b6", testResults.Suites[1].ID)

	assert.Equal(t, "foo\\zup", testResults.Suites[2].Name)
	assert.Equal(t, "5e761871-975e-34f2-ac3b-88c01a0befb2", testResults.Suites[2].ID)

	require.Len(t, testResults.Suites[0].Tests, 1)
	assert.Equal(t, "1321beac-9348-371f-a546-464e7d56304b", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)

	require.Len(t, testResults.Suites[1].Tests, 1)
	assert.Equal(t, "bd4fab01-8256-314b-861c-511a9c998c7b", testResults.Suites[1].Tests[0].ID)
	assert.Equal(t, "baz", testResults.Suites[1].Tests[0].Name)

	require.Len(t, testResults.Suites[2].Tests, 1)
	assert.Equal(t, "c77988eb-a5a0-3a27-a18e-c20f161d09ae", testResults.Suites[2].Tests[0].ID)
	assert.Equal(t, "bar", testResults.Suites[2].Tests[0].Name)

	require.Len(t, testResults.Suites[3].Tests, 2)
	assert.Equal(t, "b6215e9b-fc90-3383-a0b9-d2662094c9b2", testResults.Suites[3].Tests[0].ID)
	assert.Equal(t, "[1] testVirtualMetrics=true", testResults.Suites[3].Tests[0].Name)
	assert.Equal(t, "efe418a7-61b8-36f8-8073-05353796dc05", testResults.Suites[3].Tests[1].ID)
	assert.Equal(t, "[2] testVirtualMetrics=false", testResults.Suites[3].Tests[1].Name)

	assert.Equal(t, parser.StateError, testResults.Suites[3].Tests[0].State)
	assert.Contains(t, testResults.Suites[3].Tests[0].Error.Body, "expected: <true> but was: <false>")

}

func Test_Embedded_ParseInvalidRoot(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
		<nontestsuites name="em">
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testsuite name="zap" id="4321">
					<testcase name="baz">
					</testcase>
				</testsuite>
				<testsuite name="zup" id="54321">
					<testcase name="bar">
					</testcase>
				</testsuite>
			</testsuite>
		</nontestsuites>
	`))

	path := fileloader.Ensure(reader)

	p := NewEmbedded()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
