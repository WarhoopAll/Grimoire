package ctrl

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"grimoire/app/model"
	"strconv"
)

func (ctr *AccountHandler) AddComment(ctx *fiber.Ctx) error {
	id, ok := ctx.Locals("id").(int)
	if !ok {
		return ErrResponse(ctx, MsgUnauthorized)
	}

	entry := &model.Comment{}
	if err := ctx.BodyParser(&entry); err != nil {
		return ErrResponse(ctx, MsgInternal)
	}

	entry.Author = id

	fmt.Println(entry.Author)
	res, err := ctr.services.Sait.AddComment(ctx.Context(), entry)
	if err != nil {
		return ErrResponse(ctx, MsgInternal)
	}

	return Response(ctx, MsgSuccess, res)
}

func (ctr *AccountHandler) DeleteComment(ctx *fiber.Ctx) error {
	idacc, ok := ctx.Locals("id").(int)
	if !ok {
		return ErrResponse(ctx, MsgUnauthorized)
	}

	res, err := ctr.services.Account.GetByID(ctx.Context(), idacc)
	if err != nil {
		return ErrResponse(ctx, MsgInternal)
	}

	if res.Access == nil {
		res.Access = &model.Access{SecurityLevel: 0}
	}

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ErrResponse(ctx, MsgInternal)
	}

	comment, err := ctr.services.Sait.GetCommentByID(ctx.Context(), id)
	if err != nil {
		return ErrResponse(ctx, MsgNotFound)
	}
	if comment.Profile.AccountID != idacc && res.Access.SecurityLevel <= 0 {
		return ErrResponse(ctx, MsgForbidden)
	}

	err = ctr.services.Sait.DeleteComment(ctx.Context(), id)
	if err != nil {
		return ErrResponse(ctx, MsgInternal)
	}

	return Response(ctx, MsgSuccess, nil)
}

func (ctr *AccountHandler) UpdateComment(ctx *fiber.Ctx) error {
	idacc, ok := ctx.Locals("id").(int)
	if !ok {
		return ErrResponse(ctx, MsgUnauthorized)
	}
	entry := &model.Comment{}
	err := ctx.BodyParser(&entry)
	if err != nil {
		return ErrResponse(ctx, MsgInternal)
	}

	comment, err := ctr.services.Sait.GetCommentByID(ctx.Context(), entry.ID)
	if err != nil {
		return ErrResponse(ctx, MsgNotFound)
	}

	if comment.Profile.AccountID != idacc {
		return ErrResponse(ctx, MsgForbidden)
	}

	err = ctr.services.Sait.UpdateComment(ctx.Context(), entry)
	if err != nil {
		return ErrResponse(ctx, MsgInternal)
	}
	return Response(ctx, MsgSuccess, nil)
}
