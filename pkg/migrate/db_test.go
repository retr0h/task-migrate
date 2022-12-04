package migrate

import (
	"math/rand"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func randomName(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

var _ = Describe("Db", func() {
	var (
		logger           *logrus.Logger
		migrationName    string
		migrationRecords []string
		m                *Migration
	)

	Context("Live Database", func() {
		BeforeEach(func() {
			logger = &logrus.Logger{}
			migrationName = randomName(8)

			db, err := OpenDB()
			Expect(err).ToNot(HaveOccurred())

			repository := NewRepository(db)
			directory := "versions"
			verbose := false
			m = NewMigration(logger, repository, directory, verbose)
		})

		Context("Insert", func() {
			BeforeEach(func() {
			})

			AfterEach(func() {
				err := m.repo.Delete(migrationName)
				Expect(err).ToNot(HaveOccurred())
			})

			When("migration record does not exist", func() {
				It("should insert the given record", func() {
					err := m.repo.Insert(migrationName)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("migration record exists", func() {
				It("should return an error", func() {
					err := m.repo.Insert(migrationName)
					Expect(err).ToNot(HaveOccurred())
					err = m.repo.Insert(migrationName)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("Delete", func() {
			BeforeEach(func() {
				err := m.repo.Insert(migrationName)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				err := m.repo.Delete(migrationName)
				Expect(err).ToNot(HaveOccurred())
			})

			When("migration record exists", func() {
				It("should delete the given record", func() {
					err := m.repo.Delete(migrationName)
					Expect(err).ToNot(HaveOccurred())

					m, _ := m.repo.FindByName(migrationName)
					Expect(m).To(BeNil())
				})
			})

			When("migration record does not exist", func() {
				It("should return no error", func() {
					err := m.repo.Delete("does-not-exist")
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("FindByName", func() {
			BeforeEach(func() {
				err := m.repo.Insert(migrationName)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				err := m.repo.Delete(migrationName)
				Expect(err).ToNot(HaveOccurred())
			})

			When("migration record exists", func() {
				It("should find the given record", func() {
					m, err := m.repo.FindByName(migrationName)
					Expect(err).ToNot(HaveOccurred())
					Expect(m.Name).To(Equal(migrationName))
				})
			})

			When("migration record does not exist", func() {
				It("should return an error", func() {
					_, err := m.repo.FindByName("does-not-exist")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("Find", func() {
			BeforeEach(func() {
				migrationRecords = []string{
					randomName(8),
					randomName(8),
					randomName(8),
				}
			})

			AfterEach(func() {
				for _, record := range migrationRecords {
					err := m.repo.Delete(record)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			When("migration records exist", func() {
				It("should find all records", func() {
					for _, record := range migrationRecords {
						err := m.repo.Insert(record)
						Expect(err).ToNot(HaveOccurred())
					}

					m, err := m.repo.Find()
					Expect(err).ToNot(HaveOccurred())
					Expect(len(m)).To(Equal(3))
				})
			})

			When("migration records do not exist", func() {
				It("should return an empty slice", func() {
					m, err := m.repo.Find()
					Expect(err).ToNot(HaveOccurred())
					Expect(len(m)).To(Equal(0))
				})
			})
		})

		Context("First", func() {
			BeforeEach(func() {
				migrationRecords = []string{
					randomName(8),
					randomName(8),
					randomName(8),
				}
			})

			AfterEach(func() {
				for _, record := range migrationRecords {
					err := m.repo.Delete(record)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			When("migration record exists", func() {
				It("should find the first record", func() {
					for _, record := range migrationRecords {
						err := m.repo.Insert(record)
						Expect(err).ToNot(HaveOccurred())
					}

					m, err := m.repo.First()
					Expect(err).ToNot(HaveOccurred())
					Expect(m.Name).To(Equal(migrationRecords[0]))
				})
			})

			When("migration database empty", func() {
				It("should return nil", func() {
					m, err := m.repo.First()
					Expect(err).To(HaveOccurred())
					Expect(m).To(BeNil())
				})
			})
		})

		Context("Last", func() {
			BeforeEach(func() {
				migrationRecords = []string{
					randomName(8),
					randomName(8),
					randomName(8),
				}
			})

			AfterEach(func() {
				for _, record := range migrationRecords {
					err := m.repo.Delete(record)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			When("migration record exists", func() {
				It("should find the last record", func() {
					for _, record := range migrationRecords {
						err := m.repo.Insert(record)
						Expect(err).ToNot(HaveOccurred())
					}

					m, err := m.repo.Last()
					Expect(err).ToNot(HaveOccurred())
					Expect(m.Name).To(Equal(migrationRecords[2]))
				})
			})

			When("migration database empty", func() {
				It("should return an error", func() {
					m, err := m.repo.Last()
					Expect(err).To(HaveOccurred())
					Expect(m).To(BeNil())
				})
			})
		})
	})
})
