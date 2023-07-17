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
	// When("we use a dump from GCP's postgres 15", func() { testBindingCommonOps("15", "gcp_pg15.sql") })
	When("we use a dump from AWS's postgres 14", func() { testBindingCommonOps("14", "aws_pg14.sql") })
	// When("we use a dump from AWS's postgres 15", func() { testBindingCommonOps("15", "aws_pg15.sql") })
	When("we use a dump from AWS's aurora postgres 14", func() { testBindingCommonOps("14", "aws_aurora_pg14.sql") })
	// When("we use a dump from AWS's aurora postgres 15", func() { testBindingCommonOps("15", "aws_aurora_pg15.sql") })
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

	It("retains tables created by a binding even after the binding has been deleted. BUT NEW BINDINGS FAIL IF THERE WASN'T ANOTHER ACTIVE BINDING. THIS IS A BUG!!!!", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")
		deleteUserWorks("someuser", "someuser", factory)

		if dumpFile == "aws_pg15.sql" || dumpFile == "aws_aurora_pg15.sql" {
			createUserWorks("otheruser", "otheruser", factory)
		} else {
			createUserFails("otheruser", "otheruser", factory, "granting table privilege to datawoner role: pq: permission denied for table table1")
		}
	})

	It("retains tables created by a binding even after the binding has been deleted", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")
		createUserWorks("otheruser", "otheruser", factory)
		deleteUserWorks("someuser", "someuser", factory)

		createUserWorks("athirduser", "athirduser", factory)
	})

	It("instantiates new bindings in the same database", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlReturns("someuser", "someuser", factory, "SELECT COUNT(1) FROM pg_catalog.pg_user WHERE usename='otheruser'", "0")
		createUserWorks("otheruser", "otheruser", factory)
		customSqlReturns("someuser", "someuser", factory, "SELECT COUNT(1) FROM pg_catalog.pg_user WHERE usename='otheruser'", "1")
	})

	It("connects new bindings to public schema by default", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlReturns("someuser", "someuser", factory, "SELECT current_schema();", "public")
		createUserWorks("otheruser", "otheruser", factory)
		customSqlReturns("otheruser", "otheruser", factory, "SELECT current_schema();", "public")
	})

	It("allows any binding to create new tables", func() {
		createUserWorks("someuser", "someuser", factory)
		createUserWorks("otheruser", "otheruser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")
		customSqlFails("otheruser", "otheruser", factory, "CREATE TABLE TABLE1();", `relation "table1" already exists`)
		customSqlWorks("otheruser", "otheruser", factory, "CREATE TABLE TABLE2();")
	})

	It("prevents dropping tables owned by other existing bindings", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")

		createUserWorks("otheruser", "otheruser", factory)
		customSqlWorks("otheruser", "otheruser", factory, "CREATE TABLE TABLE2();")

		customSqlFails("someuser", "someuser", factory, "DROP TABLE TABLE2;", `pq: must be owner of table table2`)
		customSqlFails("otheruser", "otheruser", factory, "DROP TABLE TABLE1;", `pq: must be owner of table table1`)

		customSqlWorks("someuser", "someuser", factory, "DROP TABLE TABLE1;")
		customSqlWorks("otheruser", "otheruser", factory, "DROP TABLE TABLE2;")
	})

	It("allows dropping tables created by another binding only after such binding gets removed", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")

		createUserWorks("otheruser", "otheruser", factory)
		customSqlFails("otheruser", "otheruser", factory, "DROP TABLE TABLE1;", `pq: must be owner of table table1`)

		deleteUserWorks("someuser", "someuser", factory)
		customSqlWorks("otheruser", "otheruser", factory, "DROP TABLE TABLE1;")
	})

	It("prevents altering tables owned by other existing bindings", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")

		createUserWorks("otheruser", "otheruser", factory)
		customSqlWorks("otheruser", "otheruser", factory, "CREATE TABLE TABLE2();")

		customSqlFails("someuser", "someuser", factory, "ALTER TABLE TABLE2 ADD COLUMN new_column TEXT;", `pq: must be owner of table table2`)
		customSqlFails("otheruser", "otheruser", factory, "ALTER TABLE TABLE1 ADD COLUMN new_column TEXT;", `pq: must be owner of table table1`)

		customSqlWorks("someuser", "someuser", factory, "ALTER TABLE TABLE1 ADD COLUMN new_column TEXT;")
		customSqlWorks("otheruser", "otheruser", factory, "ALTER TABLE TABLE2 ADD COLUMN new_column TEXT;")
	})

	It("allows altering tables created by another binding only after such binding gets removed", func() {
		createUserWorks("someuser", "someuser", factory)
		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")

		createUserWorks("otheruser", "otheruser", factory)
		customSqlFails("otheruser", "otheruser", factory, "ALTER TABLE TABLE1 ADD COLUMN new_column TEXT;", `pq: must be owner of table table1`)

		deleteUserWorks("someuser", "someuser", factory)
		customSqlWorks("otheruser", "otheruser", factory, "ALTER TABLE TABLE1 ADD COLUMN new_column TEXT;")

		createUserWorks("athirduser", "athirduser", factory)
		customSqlWorks("athirduser", "athirduser", factory, "ALTER TABLE TABLE1 ADD COLUMN another_column TEXT;")
	})

	It("doesn't make tables visible by everyone immediately after their creation", func() {
		createUserWorks("someuser", "someuser", factory)
		createUserWorks("otheruser", "otheruser", factory)

		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")
		customSqlWorks("someuser", "someuser", factory, "SELECT COUNT(1) FROM TABLE1;")
		customSqlFails("otheruser", "otheruser", factory, "SELECT COUNT(1) FROM TABLE1;", `pq: permission denied for table table1`)
	})

	It("only makes tables visible by everyone after a new binding is created", func() {
		createUserWorks("someuser", "someuser", factory)
		createUserWorks("otheruser", "otheruser", factory)

		customSqlWorks("someuser", "someuser", factory, "CREATE TABLE TABLE1();")
		customSqlWorks("someuser", "someuser", factory, "SELECT COUNT(1) FROM TABLE1;")
		customSqlFails("otheruser", "otheruser", factory, "SELECT COUNT(1) FROM TABLE1;", `pq: permission denied for table table1`)

		createUserWorks("athirduser", "athirduser", factory)
		customSqlWorks("someuser", "someuser", factory, "SELECT COUNT(1) FROM TABLE1;")
		customSqlWorks("otheruser", "otheruser", factory, "SELECT COUNT(1) FROM TABLE1;")
		customSqlWorks("athirduser", "athirduser", factory, "SELECT COUNT(1) FROM TABLE1;")
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
			_, err = db.Exec("SELECT * FROM PUBLIC.NONEXISTING;")
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

func createUserWorks(user, password string, factory connectionFactory) {
	diag := sqlUserCreate(context.TODO(), user, password, factory)
	Expect(diag).To(BeNil())
}

func createUserFails(user, password string, factory connectionFactory, expected string) {
	diag := sqlUserCreate(context.TODO(), user, password, factory)
	Expect(diag).NotTo(BeNil())
	Expect(diag[0].Summary).To(ContainSubstring(expected))
}

func deleteUserWorks(user, password string, factory connectionFactory) {
	diag := sqlUserDelete(context.TODO(), user, password, factory)
	Expect(diag).To(BeNil())
}

func customSqlWorks(user, password string, factory connectionFactory, sql string) {
	db, err := factory.ConnectAsUser(user, password)
	Expect(err).NotTo(HaveOccurred())
	defer db.Close()
	_, err = db.Exec(sql)
	Expect(err).NotTo(HaveOccurred())
}

func customSqlFails(user, password string, factory connectionFactory, sql, expected string) {
	db, err := factory.ConnectAsUser(user, password)
	Expect(err).NotTo(HaveOccurred())
	defer db.Close()
	_, err = db.Exec(sql)
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(ContainSubstring(expected))
}

func customSqlReturns(user, password string, factory connectionFactory, sql, expected string) {
	db, err := factory.ConnectAsUser(user, password)
	Expect(err).NotTo(HaveOccurred())
	defer db.Close()
	row := db.QueryRow(sql)
	Expect(err).NotTo(HaveOccurred())

	var output string
	err = row.Scan(&output)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(Equal(expected))
}
