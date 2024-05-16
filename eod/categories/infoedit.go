package categories

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (c *Categories) ImageCmd(ctx sevcord.Ctx, cat string, image string) {
	ctx.Acknowledge()

	// Check the element
	var name string
	var old string
	err := c.db.QueryRow("SELECT name, image FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(cat), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		c.base.Error(ctx, err, "Category **"+cat+"** doesn't exist!")
		return
	}

	// Make a poll
	res := c.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindCatImage,
		Data: types.PgData{
			"cat": name,
			"new": image,
			"old": old,
		},
	})
	if !res.Ok {
		ctx.Respond(res.Response())
		return
	}

	// Create a Response
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested an image for category **%s** 📷", name)))
}

func (c *Categories) SignCmd(ctx sevcord.Ctx, opts []any) {
	// Check the element
	var name string
	var old string
	err := c.db.QueryRow("SELECT name, comment FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		c.base.Error(ctx, err, "Category **"+opts[0].(string)+"** doesn't exist!")
		return
	}

	// Get a mark
	ctx.(*sevcord.InteractionCtx).Modal(sevcord.NewModal("Sign Category", func(ctx sevcord.Ctx, s []string) {
		// Make poll
		res := c.polls.CreatePoll(ctx, &types.Poll{
			Kind: types.PollKindCatComment,
			Data: types.PgData{
				"cat": name,
				"new": s[0],
				"old": old,
			},
		})
		if !res.Ok {
			ctx.Respond(res.Response())
			return
		}

		// Respond
		ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested a note for category **%s** 🖋️", name)))
	}).Input(sevcord.NewModalInput("New Comment", "None", sevcord.ModalInputStyleParagraph, 2400)))
}

func (c *Categories) ColorCmd(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Check hex code
	code := opts[1].(string)
	if !strings.HasPrefix(code, "#") {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}
	val, err := strconv.ParseInt(strings.TrimPrefix(code, "#"), 16, 64)
	if err != nil {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}
	if val < 0 || val > 16777215 {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}

	// Check element
	var name string
	var old int
	err = c.db.QueryRow("SELECT name, color FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		c.base.Error(ctx, err, "Category **"+opts[0].(string)+"** doesn't exist!")
		return
	}

	// Make poll
	res := c.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindCatColor,
		Data: types.PgData{
			"cat": name,
			"new": float64(val),
			"old": float64(old),
		},
	})
	if !res.Ok {
		ctx.Respond(res.Response())
		return
	}

	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested a color for category **%s** 🎨", name)))
}
