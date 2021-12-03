package eodb

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (d *DB) GetElementByName(name string) (types.Element, types.GetResponse) {
	d.RLock()
	defer d.RUnlock()

	id, exists := d.elemNames[strings.ToLower(name)]
	if !exists {
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf("Element **%s** doesn't exist!", name),
		}
	}
	return d.elements[id], types.GetResponse{Exists: true}
}

func (d *DB) GetIDByName(name string) (int, types.GetResponse) {
	d.RLock()
	defer d.RUnlock()

	id, exists := d.elemNames[strings.ToLower(name)]
	if !exists {
		return 0, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf("Element **%s** doesn't exist!", name),
		}
	}
	return id, types.GetResponse{Exists: true}
}

func (d *DB) GetElement(id int) (types.Element, types.GetResponse) {
	d.RLock()
	defer d.RUnlock()
	if id < 1 {
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: "Element ID can't be negative!",
		}
	}
	if id > len(d.elements) {
		return types.Element{}, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf("Element **#%d** doesn't exist!", id),
		}
	}

	return d.elements[id-1], types.GetResponse{Exists: true}
}

func (d *DB) GetCombo(elems []int) (int, types.GetResponse) {
	txt := util.FormatCombo(elems)
	d.RLock()
	res, exists := d.combos[txt]
	d.RUnlock()
	if !exists {
		return 0, types.GetResponse{
			Exists:  false,
			Message: "Combo doesn't exist!",
		}
	}
	return res, types.GetResponse{Exists: true}
}

func (d *DB) GetInv(id string) *types.ElemContainer {
	d.RLock()
	inv, exists := d.invs[id]
	d.RUnlock()
	if !exists {
		d.Lock()
		d.invs[id] = inv
		d.Unlock()

		return types.NewElemContainer(map[int]types.Empty{
			1: {},
			2: {},
			3: {},
			4: {},
		}, id)
	}
	return inv
}

func (d *DB) GetCat(name string) (*types.ElemContainer, types.GetResponse) {
	d.RLock()
	cat, exists := d.cats[strings.ToLower(name)]
	d.RUnlock()
	if !exists {
		return nil, types.GetResponse{
			Exists:  false,
			Message: fmt.Sprintf("Category **%s** doesn't exist!", name),
		}
	}
	return cat, types.GetResponse{Exists: true}
}