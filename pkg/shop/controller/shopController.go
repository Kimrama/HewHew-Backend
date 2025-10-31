package controller

import "github.com/gofiber/fiber/v2"

type ShopController interface {
	CreateCanteen(ctx *fiber.Ctx) error
	EditCanteen(ctx *fiber.Ctx) error
	DeleteCanteen(ctx *fiber.Ctx) error
	EditShop(ctx *fiber.Ctx) error
	GetShop(ctx *fiber.Ctx) error
	GetShopByID(ctx *fiber.Ctx) error
	ChangeState(ctx *fiber.Ctx) error
	EditShopImage(ctx *fiber.Ctx) error
	GetAllCanteens(ctx *fiber.Ctx) error
	GetCanteenByName(ctx *fiber.Ctx) error
	GetAllShops(ctx *fiber.Ctx) error
	Createtag(ctx *fiber.Ctx) error
	GetTagsByShopIDAndTopic(ctx *fiber.Ctx) error
	Edittag(ctx *fiber.Ctx) error
	GetAllTags(ctx *fiber.Ctx) error
	DeleteTag(ctx *fiber.Ctx) error
	GetAllMenus(ctx *fiber.Ctx) error
	CreateTransactionLog(ctx *fiber.Ctx) error
	CreateNotification(ctx *fiber.Ctx) error
}
