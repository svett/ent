package integration_test

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/phogolabs/ent/integration/ent"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pagination", func() {
	var (
		ctx    = context.TODO()
		client *ent.Client
	)

	BeforeEach(func() {
		var err error

		// client, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", ent.Debug())
		client, err = ent.Open("postgres", "postgres://localhost:5432/ent?sslmode=disable", ent.Debug())
		Expect(err).NotTo(HaveOccurred())
		Expect(client.Schema.Create(ctx)).To(Succeed())
	})

	AfterEach(func() {
		Expect(client.Close()).To(Succeed())
	})

	Describe("Query", func() {
		var entities []*ent.Product

		BeforeEach(func() {
			entities = []*ent.Product{}

			create := func(name string) {
				i := len(entities)

				entity, err := client.Product.Create().
					SetID(imap[i]).
					SetTitle(name).
					Save(ctx)

				Expect(err).NotTo(HaveOccurred())
				entities = append(entities, entity)
			}

			create("Hat")
			create("Pants")
			create("Pants")
			create("Jackets")
			create("Hat")
			create("T-Shirt")
			create("Trousers")
			create("Cap")
			create("T-Shirt")
			create("Hat")
		})

		query := func(page *ent.Pagination) []*ent.Product {
			query := client.Product.Query().Paginate(page)

			records, err := query.All(ctx)
			Expect(err).NotTo(HaveOccurred())

			return records
		}

		AfterEach(func() {
			for _, entity := range entities {
				Expect(client.Product.DeleteOneID(entity.ID).Exec(ctx)).To(Succeed())
			}
		})

		It("returns the entities page by page", func() {
			cursor, err := ent.DecodeCursor("+title,+id", "")
			Expect(err).NotTo(HaveOccurred())

			spew.Dump(cursor)

			// navigate to page
			pagination := ent.NewPagination().
				WithCursor(cursor).
				WithLimit(2)

			// fetch first page
			records := query(pagination)
			Expect(records).To(HaveLen(2))
			Expect(records[0].Title).To(Equal("Cap"))
			Expect(records[1].Title).To(Equal("Hat"))

			// fetch next page
			cursor = cursor.Next(records)
			pagination = ent.NewPagination().
				WithCursor(cursor).
				WithLimit(2)

			records = query(pagination)
			Expect(records).To(HaveLen(2))
			Expect(records[0].Title).To(Equal("Hat"))
			Expect(records[1].Title).To(Equal("Hat"))

			// fetch next page
			cursor = cursor.Next(records)
			pagination = ent.NewPagination().
				WithCursor(cursor).
				WithLimit(2)

			records = query(pagination)
			Expect(records).To(HaveLen(2))
			Expect(records[0].Title).To(Equal("Jackets"))
			Expect(records[1].Title).To(Equal("Pants"))

			// fetch next page
			cursor = cursor.Next(records)
			pagination = ent.NewPagination().
				WithCursor(cursor).
				WithLimit(2)

			records = query(pagination)
			Expect(records).To(HaveLen(2))
			Expect(records[0].Title).To(Equal("Pants"))
			Expect(records[1].Title).To(Equal("T-Shirt"))

			// fetch next page
			cursor = cursor.Next(records)
			pagination = ent.NewPagination().
				WithCursor(cursor).
				WithLimit(2)

			records = query(pagination)
			Expect(records).To(HaveLen(2))
			Expect(records[0].Title).To(Equal("T-Shirt"))
			Expect(records[1].Title).To(Equal("Trousers"))

			cursor = cursor.Next(records)
			pagination = ent.NewPagination().
				WithCursor(cursor).
				WithLimit(2)

			records = query(pagination)
			Expect(records).To(HaveLen(0))

			cursor = cursor.Next(records)
		})
	})
})
