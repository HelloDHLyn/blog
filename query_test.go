package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testQuery(queryName, query string, args ...interface{}) (data map[string]interface{}) {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: fmt.Sprintf(query, args...),
		Context:       mockContext,
	})

	if result.HasErrors() {
		fmt.Println(result.Errors)
		Fail("Failed to execute query.")
	}

	queryData := result.Data.(map[string]interface{})[queryName]
	if queryData == nil {
		return nil
	}
	return queryData.(map[string]interface{})
}

var _ = Describe("Query", func() {
	Describe("post", func() {
		testTitle := "Awesome post 😎"
		testBody := "This is my awesome post."
		testDescription := "This is my awesome description."

		var post Post

		BeforeEach(func() {
			post = Post{
				Title:       testTitle,
				Body:        testBody,
				Description: testDescription,
				IsPublic:    true,
			}
			db.Save(&post)
		})

		It("get a post should success", func() {
			data := testQuery("post", `
			query {
				post(id: %d) {
					title
					body
					description
				}
			}`, post.ID)

			Expect(data).NotTo(BeNil())
			Expect(data["title"].(string)).To(Equal(testTitle))
			Expect(data["body"].(string)).To(Equal(testBody))
			Expect(data["description"].(string)).To(Equal(testDescription))
		})

		It("get a private post should return nil", func() {
			post.IsPublic = false
			db.Save(&post)

			data := testQuery("post", `
			query {
				post(id: %d) {
					title
					body
					description
				}
			}`, post.ID)

			Expect(data).To(BeNil())
		})

		It("get a post with invalid id should fail", func() {
			data := testQuery("post", `
			query {
				post(id : %d) {
					title
				}
			}
			`, -1)

			Expect(data).To(BeNil())
		})
	})

	Describe("postList", func() {
		testTitle := "Awesome post 😎"
		testBody := "This is my awesome post."
		testDescription := "This is my awesome description."

		var posts []Post

		BeforeEach(func() {
			for range []int{0, 1, 2, 3, 4} {
				p := Post{
					Title:       testTitle,
					Body:        testBody,
					Description: testDescription,
				}

				db.Save(&p)
				posts = append(posts, p)
			}
		})

		It("get post list should success", func() {
			data := testQuery("postList", `
			query {
				postList(page: {count: 10}) {
					items { title tagList { name } }
					pageInfo { hasNext, hasBefore }
				}
			}`)

			Expect(len(data["items"].([]interface{})) > 0).To(BeTrue())
		})
	})
})

var _ = Describe("Snippet", func() {
	Describe("snippet", func() {
		var snippet Snippet

		BeforeEach(func() {
			id, _ := uuid.NewUUID()
			snippet = Snippet{
				Title: id.String(),
				Body:  "This is my awesome snippet.",
			}
			db.Save(&snippet)
		})

		It("get a snippet by id should success", func() {
			data := testQuery("snippet", `
			query {
				snippet(id: %d) {
					title
					body
				}
			}`, snippet.ID)

			Expect(data["title"].(string)).To(Equal(snippet.Title))
			Expect(data["body"].(string)).To(Equal(snippet.Body))
		})

		It("get a snippet by title should success", func() {
			data := testQuery("snippet", `
			query {
				snippet(title: "%s") {
					title
					body
				}
			}`, snippet.Title)

			Expect(data["title"].(string)).To(Equal(snippet.Title))
			Expect(data["body"].(string)).To(Equal(snippet.Body))
		})

		It("get a snippet with invalid id should return nil", func() {
			data := testQuery("snippet", `
			query {
				snippet(id : %d) {
					title
				}
			}`, -1)

			Expect(data).To(BeNil())
		})

		It("get a snippet with invalid title should  return nil", func() {
			data := testQuery("snippet", `
			query {
				snippet(title: "%s") {
					title
				}
			}`, snippet.Title+" invalid")

			Expect(data).To(BeNil())
		})
	})

	Describe("snippetList", func() {
		var snippets []Snippet

		BeforeEach(func() {
			for range []int{0, 1, 2, 3, 4} {
				id, _ := uuid.NewUUID()
				p := Snippet{
					Title: id.String(),
					Body:  "This is my awesome snippet.",
				}

				db.Save(&p)
				snippets = append(snippets, p)
			}
		})

		It("get snippet list should success", func() {
			data := testQuery("snippetList", `
			query {
				snippetList(page: {count: 10}) {
					items { title }
					pageInfo { hasNext, hasBefore }
				}
			}`)

			Expect(len(data["items"].([]interface{})) > 0).To(BeTrue())
		})
	})
})
