package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomBlock(t *testing.T, user1, user2 User) Blockuser {
	newBlock := CreateBlockParams{
		UserID:      user1.UserID,
		BlockuserID: user2.UserID,
		Reason:      gofakeit.Sentence(10),
		Status:      util.RandomStatus(),
	}

	block, err := testQueries.CreateBlock(context.Background(), newBlock)
	require.NoError(t, err)
	require.NotEmpty(t, block)
	require.Equal(t, newBlock.UserID, block.UserID)
	require.Equal(t, newBlock.BlockuserID, block.BlockuserID)
	require.Equal(t, newBlock.Reason, block.Reason)
	require.Equal(t, newBlock.Status, block.Status)
	require.NotEmpty(t, block.BlockAt)

	return block
}

func TestGetBlockUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)
	CreateRandomBlock(t, user1, user2)
}

func TestGetAllBlockUsersList(t *testing.T) {
	for i := 0; i < 20; i++ {
		user1 := CreateRandomUser(t)
		user2 := CreateRandomUser(t)
		CreateRandomBlock(t, user1, user2)
	}

	allBlockData, err := testQueries.GetAllBlockUsersList(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, allBlockData)
	require.GreaterOrEqual(t, len(allBlockData), 1)
}

func TestGetBlockUserlist(t *testing.T) {
	user1 := CreateRandomUser(t)
	for i := 0; i < 20; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomBlock(t, user1, user2)
	}

	blockList, err := testQueries.GetBlockUserlist(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, blockList)
	require.GreaterOrEqual(t, len(blockList), 1)
}

func TestUnBlockUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)
	oldBlock := CreateRandomBlock(t, user1, user2)

	var newStatus string

	for {
		newStatus = util.RandomStatus()
		if newStatus != oldBlock.Status {
			break
		}
	}

	unBlockData := UnBlockUserParams{
		UserID:      user1.UserID,
		BlockuserID: user2.UserID,
		Status:      newStatus,
	}

	unBlocked, err := testQueries.UnBlockUser(context.Background(), unBlockData)
	require.NoError(t, err)
	require.NotEmpty(t, unBlocked)
	require.NotEqual(t, unBlockData.Status, oldBlock.Status)
	require.NotEmpty(t, unBlocked.UnblockAt)
}
