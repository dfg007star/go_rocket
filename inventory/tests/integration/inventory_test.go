//go:build integration

package integration

import (
	"context"

	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции inventory")

		cancel()
	})

	Describe("GetPart", func() {
		It("должен успешно получать созданную ранее запчасть", func() {
			partUuid, _ := env.InsertTestPart(ctx)

			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
		})
	})

	Describe("ListParts All", func() {
		It("должен успешно возвращать все запчасти", func() {
			partsCount := 5
			partsUuids := make([]string, 0, partsCount)
			for i := 0; i < partsCount; i++ {
				partUuid, _ := env.InsertTestPart(ctx)
				partsUuids = append(partsUuids, partUuid)
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp.Parts)).To(Equal(len(partsUuids)))
		})
	})

	Describe("ListParts By Uuid", func() {
		It("должен успешно возвращать запчасть(и) по Uuid", func() {
			partUuid, _ := env.InsertTestPart(ctx)

			filter := &inventoryV1.PartsFilter{
				Uuids: []string{partUuid},
			}
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Parts).ToNot(BeEmpty())
			Expect(resp.Parts).To(HaveLen(1))
			Expect(resp.Parts[0].Uuid).To(Equal(partUuid))
		})
	})
})
