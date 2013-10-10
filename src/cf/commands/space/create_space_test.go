package space_test

import (
	. "cf/commands/space"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestCreateSpaceFailsWithUsage(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{}
	spaceRepo := &testapi.FakeSpaceRepository{}

	fakeUI := callCreateSpace([]string{}, reqFactory, spaceRepo)
	assert.True(t, fakeUI.FailedWithUsage)

	fakeUI = callCreateSpace([]string{"my-space"}, reqFactory, spaceRepo)
	assert.False(t, fakeUI.FailedWithUsage)
}

func TestCreateSpaceRequirements(t *testing.T) {
	spaceRepo := &testapi.FakeSpaceRepository{}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}
	callCreateSpace([]string{"my-space"}, reqFactory, spaceRepo)
	assert.True(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: false}
	callCreateSpace([]string{"my-space"}, reqFactory, spaceRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)

	reqFactory = &testreq.FakeReqFactory{LoginSuccess: false, TargetedOrgSuccess: true}
	callCreateSpace([]string{"my-space"}, reqFactory, spaceRepo)
	assert.False(t, testcmd.CommandDidPassRequirements)

}

func TestCreateSpace(t *testing.T) {
	spaceRepo := &testapi.FakeSpaceRepository{}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}
	fakeUI := callCreateSpace([]string{"my-space"}, reqFactory, spaceRepo)

	assert.Contains(t, fakeUI.Outputs[0], "Creating space")
	assert.Contains(t, fakeUI.Outputs[0], "my-space")
	assert.Equal(t, spaceRepo.CreateSpaceName, "my-space")
	assert.Contains(t, fakeUI.Outputs[1], "OK")
	assert.Contains(t, fakeUI.Outputs[2], "TIP")
}

func TestCreateSpaceWhenItAlreadyExists(t *testing.T) {
	spaceRepo := &testapi.FakeSpaceRepository{CreateSpaceExists: true}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}
	fakeUI := callCreateSpace([]string{"my-space"}, reqFactory, spaceRepo)

	assert.Equal(t, len(fakeUI.Outputs), 3)
	assert.Contains(t, fakeUI.Outputs[1], "OK")
	assert.Contains(t, fakeUI.Outputs[2], "my-space")
	assert.Contains(t, fakeUI.Outputs[2], "already exists")
}

func callCreateSpace(args []string, reqFactory *testreq.FakeReqFactory, spaceRepo *testapi.FakeSpaceRepository) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("create-space", args)

	cmd := NewCreateSpace(ui, spaceRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
