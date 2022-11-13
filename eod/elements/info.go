package elements

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
)

func (e *Elements) Info(c sevcord.Ctx, params []any) {
	c.Acknowledge()

	// Get element
	var elem types.Element
	err := e.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", params[0].(int64), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Check if you have
	description := elem.Comment
	var have bool
	err = e.db.QueryRow(`SELECT $1=ANY(inv) FROM inventories WHERE guild=$2 AND "user"=$3`, elem.ID, c.Guild(), c.Author().User.ID).Scan(&have)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	if have {
		description = "**You have this.**\n\n" + description
	}

	// Embed
	emb := sevcord.NewEmbed().
		Title(elem.Name+" Info").
		Description(description).
		Color(elem.Color).
		AddField("Creator", fmt.Sprintf("<@%s>", elem.Creator), true).
		AddField("Created On", fmt.Sprintf("<t:%d>", elem.CreatedOn.Unix()), true).
		AddField("Tree Size", util.FormatInt(elem.TreeSize), true)

	// Optional things
	if elem.Image != "" {
		emb.Thumbnail(elem.Image)
	}
	if elem.Commenter != "" {
		emb.AddField("Commenter", fmt.Sprintf("<@%s>", elem.Commenter), true)
	}
	if elem.Colorer != "" {
		emb.AddField("Colorer", fmt.Sprintf("<@%s>", elem.Colorer), true)
	}
	if elem.Imager != "" {
		emb.AddField("Imager", fmt.Sprintf("<@%s>", elem.Imager), true)
	}

	// Respond
	msg := sevcord.NewMessage("").AddEmbed(emb)
	c.Respond(msg)
}
