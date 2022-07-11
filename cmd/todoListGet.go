/*
Copyright © 2022 Remi Ferrand

Contributor(s): Remi Ferrand <riton.github_at_gmail(dot)com>, 2022

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.

*/
package cmd

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	apiV1 "github.com/toutdoux-app/toutdoux-cli/api/v1"
)

// todoListCmd represents the todoList command
var todoListGetCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: todoListGetRunE,
	Args: cobra.ExactArgs(1),
}

func init() {
	todoListCmd.AddCommand(todoListGetCmd)
}

func todoListGetRunE(cmd *cobra.Command, args []string) error {
	var todoListID uuid.UUID

	idOrName := newUUIDOrName(args[0])
	switch idOrName.Type() {
	case uuidOrNameIsUUID:
		todoListID = idOrName.uuid

	case uuidOrNameIsName:
		lists, err := apiClient.GetTodoLists(apiV1.GetTodoListsOptions{})
		if err != nil {
			return errors.Wrap(err, "listing todo lists")
		}

		list, err := lists.GetByName(idOrName.value)
		if err != nil {
			return errors.Wrapf(err, "no such todo list %s", idOrName.value)
		}

		todoListID = list.ID
	}

	todoList, err := apiClient.GetTodoListByID(todoListID.String(), apiV1.GetTodoListByIDOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("todoList = %+v\n", todoList)

	return nil
}
