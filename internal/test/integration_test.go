package test

import (
    "testing"
    "time"

    "github.com/wyw14/solo-037/internal/carriage"
    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/tank"
    "github.com/wyw14/solo-037/internal/track"
)

func TestArmOccupiedSectionNotRegistered(t *testing.T) {
    at := time.Unix(1700000000, 0)
    registry := NewRegistry()
    tracks := track.NewCoordinator(tank.NewRegistry(2.4))
    startup := Startup{Tests: registry, Tracks: tracks, Carriage: carriage.NewController()}

    first := model.NewOperation("trial-1", "hull", "section-a", "carriage-1", at)
    lease, err := startup.Arm(first, first.Section, at)
    if err != nil {
        t.Fatalf("first arm should succeed: %v", err)
    }
    if lease.Section != "section-a" {
        t.Fatalf("first arm should lease section-a, got %q", lease.Section)
    }

    second := model.NewOperation("trial-2", "hull", "section-a", "carriage-2", at)
    if _, err := startup.Arm(second, second.Section, at); err == nil {
        t.Fatal("second arm on occupied section should fail")
    }

    if _, found := registry.Get(second.ID); found {
        t.Fatal("failed application must not enter the test registry")
    }
    active := registry.Active()
    if len(active) != 1 || active[0].ID != first.ID {
        t.Fatalf("runnable list should contain only the first trial, got %+v", active)
    }
}
