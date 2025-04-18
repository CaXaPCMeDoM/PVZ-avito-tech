package entity_test

import (
	"PVZ-avito-tech/internal/entity"
	"testing"
)

func TestCity_IsValidCity(t *testing.T) {
	tests := []struct {
		name string
		city entity.City
		want bool
	}{
		{
			name: "valid city - Moscow",
			city: entity.CityMoscow,
			want: true,
		},
		{
			name: "valid city - St. Petersburg",
			city: entity.CitySpb,
			want: true,
		},
		{
			name: "valid city - Kazan",
			city: entity.CityKazan,
			want: true,
		},
		{
			name: "invalid city",
			city: "Invalid City",
			want: false,
		},
		{
			name: "empty city",
			city: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.city.IsValidCity(); got != tt.want {
				t.Errorf("City.IsValidCity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCity_ValidateCity(t *testing.T) {
	tests := []struct {
		name    string
		city    entity.City
		wantErr bool
	}{
		{
			name:    "valid city - Moscow",
			city:    entity.CityMoscow,
			wantErr: false,
		},
		{
			name:    "valid city - St. Petersburg",
			city:    entity.CitySpb,
			wantErr: false,
		},
		{
			name:    "valid city - Kazan",
			city:    entity.CityKazan,
			wantErr: false,
		},
		{
			name:    "invalid city",
			city:    "Invalid City",
			wantErr: true,
		},
		{
			name:    "empty city",
			city:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.city.ValidateCity()
			if (err != nil) != tt.wantErr {
				t.Errorf("City.ValidateCity() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				if err.Error() != "invalid city" {
					t.Errorf("City.ValidateCity() error message = %v, want %v", err.Error(), "invalid city")
				}
			}
		})
	}
}

func TestProductType_IsValidProductType(t *testing.T) {
	tests := []struct {
		name        string
		productType entity.ProductType
		want        bool
	}{
		{
			name:        "valid product type - shoes",
			productType: entity.ShoesProductType,
			want:        true,
		},
		{
			name:        "valid product type - clothes",
			productType: entity.ClothesProductType,
			want:        true,
		},
		{
			name:        "valid product type - electronics",
			productType: entity.ElectronicsProductType,
			want:        true,
		},
		{
			name:        "invalid product type",
			productType: "Invalid Product Type",
			want:        false,
		},
		{
			name:        "empty product type",
			productType: "",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.productType.IsValidProductType(); got != tt.want {
				t.Errorf("ProductType.IsValidProductType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductType_ValidateProductType(t *testing.T) {
	tests := []struct {
		name        string
		productType entity.ProductType
		wantErr     bool
	}{
		{
			name:        "valid product type - shoes",
			productType: entity.ShoesProductType,
			wantErr:     false,
		},
		{
			name:        "valid product type - clothes",
			productType: entity.ClothesProductType,
			wantErr:     false,
		},
		{
			name:        "valid product type - electronics",
			productType: entity.ElectronicsProductType,
			wantErr:     false,
		},
		{
			name:        "invalid product type",
			productType: "Invalid Product Type",
			wantErr:     true,
		},
		{
			name:        "empty product type",
			productType: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.productType.ValidateProductType()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductType.ValidateProductType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				if err.Error() != "invalid product type" {
					t.Errorf("ProductType.ValidateProductType() error message = %v, want %v", err.Error(), "invalid product type")
				}
			}
		})
	}
}

func TestReceptionsStatus_IsValidReceptionsStatus(t *testing.T) {
	tests := []struct {
		name   string
		status entity.ReceptionsStatus
		want   bool
	}{
		{
			name:   "valid status - in progress",
			status: entity.InProgressStatus,
			want:   true,
		},
		{
			name:   "valid status - close",
			status: entity.CloseStatus,
			want:   true,
		},
		{
			name:   "invalid status",
			status: "Invalid Status",
			want:   false,
		},
		{
			name:   "empty status",
			status: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValidReceptionsStatus(); got != tt.want {
				t.Errorf("ReceptionsStatus.IsValidReceptionsStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceptionsStatus_ValidateReceptionsStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  entity.ReceptionsStatus
		wantErr bool
	}{
		{
			name:    "valid status - in progress",
			status:  entity.InProgressStatus,
			wantErr: false,
		},
		{
			name:    "valid status - close",
			status:  entity.CloseStatus,
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  "Invalid Status",
			wantErr: true,
		},
		{
			name:    "empty status",
			status:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.status.ValidateReceptionsStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReceptionsStatus.ValidateReceptionsStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				if err.Error() != "invalid reception status" {
					t.Errorf("ReceptionsStatus.ValidateReceptionsStatus() error message = %v, want %v", err.Error(), "invalid reception status")
				}
			}
		})
	}
}

func TestUserRole_IsValidRole(t *testing.T) {
	tests := []struct {
		name string
		role entity.UserRole
		want bool
	}{
		{
			name: "valid role - employee",
			role: entity.UserRoleEmployee,
			want: true,
		},
		{
			name: "valid role - moderator",
			role: entity.UserRoleModerator,
			want: true,
		},
		{
			name: "invalid role",
			role: "Invalid Role",
			want: false,
		},
		{
			name: "empty role",
			role: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValidRole(); got != tt.want {
				t.Errorf("UserRole.IsValidRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRole_ValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		role    entity.UserRole
		wantErr bool
	}{
		{
			name:    "valid role - employee",
			role:    entity.UserRoleEmployee,
			wantErr: false,
		},
		{
			name:    "valid role - moderator",
			role:    entity.UserRoleModerator,
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    "Invalid Role",
			wantErr: true,
		},
		{
			name:    "empty role",
			role:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.role.ValidateRole()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRole.ValidateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				if err.Error() != "invalid user role" {
					t.Errorf("UserRole.ValidateRole() error message = %v, want %v", err.Error(), "invalid user role")
				}
			}
		})
	}
}
