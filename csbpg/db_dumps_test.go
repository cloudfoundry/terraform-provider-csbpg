package csbpg

import (
	"context"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("for every supported Postgres and IAAS combination", Label("dumps"), func() {
	When("we use a dump from GCP's postgres 14", func() { testBindingCommonOps("14", "gcp_pg14.sql") })
	When("we use a dump from GCP's postgres 15", func() { testBindingCommonOps("15", "gcp_pg15.sql") })
	When("we use a dump from AWS's postgres 14", func() { testBindingCommonOps("14", "aws_pg14.sql") })
	When("we use a dump from AWS's postgres 15", func() { testBindingCommonOps("15", "aws_pg15.sql") })
	When("we use a dump from AWS's aurora postgres 14", func() { testBindingCommonOps("14", "aws_aurora_pg14.sql") })
	When("we use a dump from AWS's aurora postgres 15", func() { testBindingCommonOps("15", "aws_aurora_pg15.sql") })
})

func testBindingCommonOps(pgVersion, dumpFile string) {
	var factory connectionFactory
	var err error

	BeforeEach(func() {
		factory, err = preparePostgresInstance(pgVersion, dumpFile)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err = cleanPostgresInstance(pgVersion, dumpFile)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can perform all common operations with bindings", func() {
		ctx := context.TODO()

		By("creating a new user", func() {
			diag := sqlUserCreate(ctx, "someuser", "someuser", factory)
			Expect(diag).To(BeNil())
		})

		By("creating a table as the new user", func() {
			db, err := factory.ConnectAsUser("someuser", "someuser")
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()
			_, err = db.Exec("CREATE TABLE PUBLIC.AAA();")
			Expect(err).NotTo(HaveOccurred())
		})

		By("creating a second user", func() {
			diag := sqlUserCreate(ctx, "otheruser", "otheruser", factory)
			Expect(diag).To(BeNil())
		})

		By("deleting the first user", func() {
			diag := sqlUserDelete(ctx, "someuser", "someuser", factory)
			Expect(diag).To(BeNil())
		})

		By("reading the table created by the now deleted first user", func() {
			db, err := factory.ConnectAsUser("otheruser", "otheruser")
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()
			_, err = db.Exec("SELECT * FROM PUBLIC.AAA;")
			Expect(err).NotTo(HaveOccurred())
		})

		By("failing to read a non existing table", func() {
			db, err := factory.ConnectAsUser("otheruser", "otheruser")
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()
			_, err = db.Exec("SELECT * FROM PUBLIC.NON-EXISTING;")
			Expect(err).To(HaveOccurred())
		})

	})
}

func preparePostgresInstance(pgVersion, dumpFile string) (connectionFactory, error) {
	cmd := exec.Command("/bin/bash", "-c", `
		docker build \
			--no-cache --tag "${IMAGE_TAG}"              \
			--build-arg PG_VERSION="${PG_VERSION}"       \
			--build-arg DUMP_FILE="${DUMP_FILE}"         \
			../testfixtures
		docker run -d --rm --name "test" -p 5999:5432 "${IMAGE_TAG}"
		until [[ "$(docker inspect -f \{\{.State.Health.Status\}\} test)" == "healthy" ]]; do
			sleep 0.1;
		done;
	`)
	cmd.Env = append(cmd.Env, "PG_VERSION="+pgVersion)
	cmd.Env = append(cmd.Env, "DUMP_FILE="+dumpFile)
	cmd.Env = append(cmd.Env, "IMAGE_TAG="+strings.Split(dumpFile, ".")[0])

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	if err != nil {
		return connectionFactory{}, err
	}
	defer session.Terminate()
	Eventually(session, 180).Should(gexec.Exit(0))

	return connectionFactory{
		host:          "localhost",
		port:          5999,
		username:      "testuser",
		password:      "password-test",
		database:      "testdb",
		dataOwnerRole: "binding_user_group",
		sslMode:       "disable",
	}, nil
}

func cleanPostgresInstance(pgVersion, dumpFile string) error {
	cmd := exec.Command("/bin/sh", "-c", `
		docker rm -f test
		docker image rm -f ${IMAGE_TAG}
		docker image rm -f postgres:${PG_VERSION}
	`)
	cmd.Env = append(cmd.Env, "PG_VERSION="+pgVersion)
	cmd.Env = append(cmd.Env, "IMAGE_TAG="+strings.Split(dumpFile, ".")[0])

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	if err != nil {
		return err
	}
	defer session.Terminate()
	Eventually(session, 120).Should(gexec.Exit())
	return nil
}
