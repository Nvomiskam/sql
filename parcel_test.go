package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)

	require.NoError(t, err)
	require.NotNil(t, id)

	testParcel, err := store.Get(id)

	require.NoError(t, err)

	assert.Equal(t, testParcel.Client, parcel.Client)
	assert.Equal(t, testParcel.Status, parcel.Status)
	assert.Equal(t, testParcel.Address, parcel.Address)
	assert.Equal(t, testParcel.CreatedAt, parcel.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)

	require.NoError(t, err)
	assert.NotNil(t, id)

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	testParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, testParcel.Address, newAddress)
}

func TestSetStatus(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)

	require.NoError(t, err)
	assert.NotNil(t, id)

	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	testParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, testParcel.Status, newStatus)
}

func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])

		require.NoError(t, err)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		expectedParcel := parcelMap[parcel.Number]
		assert.Equal(t, expectedParcel, parcel)
	}
}
