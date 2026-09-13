package settlement

import "errors"

// WorkerRole represents a type of job
type WorkerRole string

const (
	RoleFarmer     WorkerRole = "Petani"
	RoleLumberjack WorkerRole = "Penebang"
	RoleMiner      WorkerRole = "Penambang"
	RoleBlacksmith WorkerRole = "Pandai Besi"
	RoleMilitia    WorkerRole = "Garda Milisi"
)

var (
	ErrNoUnassignedSettlers = errors.New("tidak ada warga menganggur yang tersedia")
	ErrNoWorkersInRole      = errors.New("tidak ada pekerja dalam peran ini untuk ditarik")
	ErrBuildingRequired     = errors.New("butuh bangunan khusus untuk peran ini")
)

// AssignWorker increments a worker role if an unassigned settler is available
func (s *Settlement) AssignWorker(role WorkerRole) error {
	if s.UnassignedSettlers() <= 0 {
		return ErrNoUnassignedSettlers
	}

	switch role {
	case RoleFarmer:
		s.Workers.Farmers++
	case RoleLumberjack:
		s.Workers.Lumberjacks++
	case RoleMiner:
		s.Workers.Miners++
	case RoleBlacksmith:
		if s.Buildings[BuildingBlacksmith] < 1 {
			return errors.New("butuh Bengkel Pandai Besi (Level 1) untuk menugaskan Pandai Besi")
		}
		s.Workers.Blacksmiths++
	case RoleMilitia:
		s.Workers.Militia++
		s.RecalculateDefense()
	default:
		return errors.New("peran pekerja tidak valid")
	}

	return nil
}

// UnassignWorker decrements a worker role and frees up a settler
func (s *Settlement) UnassignWorker(role WorkerRole) error {
	switch role {
	case RoleFarmer:
		if s.Workers.Farmers <= 0 {
			return ErrNoWorkersInRole
		}
		s.Workers.Farmers--
	case RoleLumberjack:
		if s.Workers.Lumberjacks <= 0 {
			return ErrNoWorkersInRole
		}
		s.Workers.Lumberjacks--
	case RoleMiner:
		if s.Workers.Miners <= 0 {
			return ErrNoWorkersInRole
		}
		s.Workers.Miners--
	case RoleBlacksmith:
		if s.Workers.Blacksmiths <= 0 {
			return ErrNoWorkersInRole
		}
		s.Workers.Blacksmiths--
	case RoleMilitia:
		if s.Workers.Militia <= 0 {
			return ErrNoWorkersInRole
		}
		s.Workers.Militia--
		s.RecalculateDefense()
	default:
		return errors.New("peran pekerja tidak valid")
	}

	return nil
}
