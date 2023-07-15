package service

import (
	"context"
	"errors"
	"log"
	pb "service1/pb"
	"service1/pkg/entity"
	repo "service1/pkg/repository"
	"service1/pkg/utils"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type MyService struct {
	pb.UnimplementedMyServiceServer
}

func (s *MyService) MyMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Println("Microservice1: MyMethod called")

	result := "Hello, " + req.Data
	return &pb.Response{Result: result}, nil
}

func (s *MyService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("Microservice1: CreateUser called")
	email, err := repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email != nil {
		return nil, errors.New("user with this email already exists")
	}
	phone, err := repo.GetByPhone(req.Phone)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if phone != nil {
		return nil, errors.New("user with this phone no already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
	}

	err = repo.Create(newUser)
	if err != nil {
		return nil, err
	} else {
		result := "user created succesfuly"
		return &pb.CreateUserResponse{Firstname: req.FirstName, Email: req.Email, Result: result}, nil
	}
}

func (s *MyService) CreateUserWithOtp(ctx context.Context, req *pb.CreateUserWithOtpRequest) (*pb.CreateUserWithOtpResponse, error) {
	var otpKey entity.OtpKey
	email, err := repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email != nil {
		return nil, errors.New("user with this email already exists")
	}
	phone, err := repo.GetByPhone(req.Phone)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if phone != nil {
		return nil, errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.Password = string(hashedPassword)
	key, err := utils.SendOtp(req.Phone)
	if err != nil {
		return nil, err
	} else {
		var user entity.Signup
		copier.Copy(&user, &req)
		err = repo.CreateSignup(&user)
		otpKey.Key = key
		otpKey.Phone = req.Phone
		err = repo.CreateOtpKey(&otpKey)
		if err != nil {
			return nil, err
		}
		result := "Otp send succesfuly"
		return &pb.CreateUserWithOtpResponse{Phone: req.Phone, Key: key, Result: result}, nil
	}
}

func (s *MyService) SignupOtpValidation(ctx context.Context, req *pb.OtpValidationRequest) (*pb.OtpValidationResponse, error) {
	result, err := repo.GetByKey(req.Key)
	if err != nil {
		return nil, err
	}
	user, err := repo.GetSignupByPhone(result.Phone)
	if err != nil {
		return nil, err
	}
	err = utils.CheckOtp(result.Phone, req.Otp)
	if err != nil {
		return nil, err
	} else {
		newUser := &entity.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Password:  user.Password,
		}

		err = repo.Create(newUser)
		if err != nil {
			return nil, err
		} else {
			result := "Otp send succesfuly"
			return &pb.OtpValidationResponse{Result: result}, nil
		}
	}
}

func (s *MyService) LoginWithOtp(ctx context.Context, req *pb.LoginWithOtpRequest) (*pb.LoginWithOtpResponse, error) {
	var otpKey entity.OtpKey
	result, err := repo.GetByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("user with this phone not found")
	}
	permission, err := repo.CheckPermission(result)
	if permission == false {
		return nil, errors.New("user permission denied")
	}
	key, err := utils.SendOtp(req.Phone)
	if err != nil {
		return nil, err
	} else {
		otpKey.Key = key
		otpKey.Phone = req.Phone
		err = repo.CreateOtpKey(&otpKey)
		if err != nil {
			return nil, err
		}
		result := "Otp send succesfuly"
		return &pb.LoginWithOtpResponse{Key: key, Result: result, Phone: req.Phone}, nil
	}
}

func (s *MyService) LoginOtpValidation(ctx context.Context, req *pb.OtpValidationRequest) (*pb.LoginOtpValidationResponse, error) {
	result, err := repo.GetByKey(req.Key)
	if err != nil {
		return nil, err
	}
	user, err := repo.GetByPhone(result.Phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(result.Phone, req.Otp)
	if err1 != nil {
		return nil, err1
	} else {
		result := "User Loged in succesfuly"
		return &pb.LoginOtpValidationResponse{Userid: int32(user.ID), Result: result}, nil
	}

}

func (s *MyService) LoginWithPassword(ctx context.Context, req *pb.LoginWithPasswordRequest) (*pb.LoginWithPasswordResponse, error) {
	user, err := repo.GetByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user with this phone not found")
	}
	permission, err := repo.CheckPermission(user)
	if permission == false {
		return nil, errors.New("user permission denied")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Invalid Password")
	} else {
		id := int32(user.ID)
		result := "user loged in succesfuly and cookie stored"
		return &pb.LoginWithPasswordResponse{Userid: id, Result: result}, nil
	}
}

func (s *MyService) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest) (*pb.RegisterAdminResponse, error) {
	email, err := repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if email != nil {
		return nil, errors.New("admin with this email already exists")
	}
	phone, err := repo.GetByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if phone != nil {
		return nil, errors.New("admin with this phone already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newAdmin := &entity.Admin{
		AdminName: req.Adminname,
		Email:     req.Email,
		Phone:     req.Phone,
		Role:      req.Role,
		Password:  string(hashedPassword),
	}

	err = repo.AdminCreate(newAdmin)
	if err != nil {
		return nil, err
	} else {

	}
	result := "user loged in succesfuly and cookie stored"
	return &pb.RegisterAdminResponse{Result: result}, nil
}

func (s *MyService) AdminLoginWithPassword(ctx context.Context, req *pb.LoginWithPasswordRequest) (*pb.LoginWithPasswordResponse, error) {
	admin, err := repo.GetByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("admin with this phone not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Invalid Password")
	} else {
		result := "user loged in succesfuly and cookie stored"
		return &pb.LoginWithPasswordResponse{Userid: int32(admin.ID), Result: result}, nil
	}
}

// oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo

// func (s *MyService) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.Response, error) {

// 	newAddress := &entity.Address{
// 		UserId:  int(req.Userid),
// 		House:   req.House,
// 		Street:  req.Street,
// 		City:    req.City,
// 		Pincode: req.Pincode,
// 		Type:    req.Type,
// 	}
// 	err := repo.CreateAddress(newAddress)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		result := "user address added succesfuly"
// 		return &pb.Response{Result: result}, nil
// 	}

// }

// func (s *MyService) AdminSignup(ctx context.Context, req *pb.AdminSignupRequest) (*pb.Response, error) {
// 	email, err := repo.AdminGetByEmail(req.Email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if email != nil {
// 		return nil, errors.New("admin with this email already exists")
// 	}
// 	phone, err := repo.AdminGetByPhone(req.Phone)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if phone != nil {
// 		return nil, errors.New("admin with this phone already exists")
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	newAdmin := &entity.Admin{
// 		AdminName: req.Adminname,
// 		Email:     req.Email,
// 		Phone:     req.Phone,
// 		Role:      req.Role,
// 		Password:  string(hashedPassword),
// 	}

// 	err = repo.AdminCreate(newAdmin)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := "admin created succesfuly"
// 	return &pb.Response{Result: result}, nil
// }

// func (s *MyService) AdminLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
// 	user, err := repo.GetByPhone(req.Phone)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		return nil, errors.New("user with this phone not found")
// 	}
// 	permission, err := repo.CheckPermission(user)
// 	if permission == false {
// 		return nil, errors.New("user permission denied")
// 	}
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
// 		return nil, errors.New("Invalid Password")
// 	} else {
// 		id := int32(user.ID)
// 		result := "admin loged in succesfuly and cookie stored"
// 		return &pb.LoginResponse{Userid: id, Result: result}, nil
// 	}
// }
