package repository

import (
	"testing"
)

func TestColumnIterator1(t *testing.T) {
	iter := columnIterator{columns: []string{"name", "email"}, identity: "id", i: -1}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "name" {
		t.Error("first <> 'name'")
	}
	if iter.Get() != "name" {
		t.Error("first <> 'name'")
	}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "email" {
		t.Error("first <> 'email'")
	}
	if iter.Get() != "email" {
		t.Error("first <> 'email'")
	}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "id" {
		t.Error("first <> 'id'")
	}
	if iter.Get() != "id" {
		t.Error("first <> 'id'")
	}
	if iter.Next() {
		t.Error("has next")
	}
}

func TestColumnIterator2(t *testing.T) {
	iter := columnIterator{columns: []string{"id", "name", "email"}, identity: "id", i: -1}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "id" {
		t.Error("first <> 'id'")
	}
	if iter.Get() != "id" {
		t.Error("first <> 'id'")
	}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "name" {
		t.Error("first <> 'name'")
	}
	if iter.Get() != "name" {
		t.Error("first <> 'name'")
	}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "email" {
		t.Error("first <> 'email'")
	}
	if iter.Get() != "email" {
		t.Error("first <> 'email'")
	}
	if iter.Next() {
		t.Error("has next")
	}
}

func TestColumnIterator3(t *testing.T) {
	iter := columnIterator{columns: []string{"name", "email", "id"}, identity: "id", i: -1}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "name" {
		t.Error("first <> 'name'")
	}
	if iter.Get() != "name" {
		t.Error("first <> 'name'")
	}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "email" {
		t.Error("first <> 'email'")
	}
	if iter.Get() != "email" {
		t.Error("first <> 'email'")
	}
	if !iter.Next() {
		t.Error("no next")
	}
	if iter.Get() != "id" {
		t.Error("first <> 'id'")
	}
	if iter.Get() != "id" {
		t.Error("first <> 'id'")
	}
	if iter.Next() {
		t.Error("has next")
	}
}
