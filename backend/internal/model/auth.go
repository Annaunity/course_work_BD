package model

type User struct {
  Id       string `json:"id"`
  FullName string `json:"full_name"`
  Street   string `json:"street"`
  House    string `json:"house"`
  Email    string `json:"email"`
  Phone    string `json:"phone"`
  Password string `json:"password"`
  IsAdmin  string `json:"is_admin"`
}

type CreateUserDTO struct {
  FullName string `json:"full_name"`
  Street   string `json:"street"`
  House    string `json:"house"`
  Password string `json:"password"`
}

type UpdateUserDTO struct {
  FullName string `json:"full_name"`
  Street   string `json:"street"`
  House    string `json:"house"`
  Email    string `json:"email"`
  Phone    string `json:"phone"`
}

type AuthUserDTO struct {
  FullName string `json:"full_name"`
  Password string `json:"password"`
}

func (dto *CreateUserDTO) NewUser() *User {
  return &User{
    FullName: dto.FullName,
    Street:   dto.Street,
    House:    dto.House,
    Password: dto.Password,
  }
}

func (dto *AuthUserDTO) NewUser() *User {
  return &User{
    FullName: dto.FullName,
    Password: dto.Password,
  }
}

func (dto *UpdateUserDTO) NewUser() *User {
  return &User{
    FullName: dto.FullName,
    Street:   dto.Street,
    House:    dto.House,
    Email:    dto.Email,
    Phone:    dto.Phone,
  }
}
