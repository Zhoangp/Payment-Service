package client

import (
	"fmt"
	"github.com/Zhoangp/Payment-Service/config"
	"github.com/Zhoangp/Payment-Service/pb"
	"github.com/Zhoangp/Payment-Service/pb/course"
	"google.golang.org/grpc"
)

func InitCartServiceClient(c *config.Config) (pb.CartServiceClient, error) {
	// using WithInsecure() because no SSL running

	cc, err := grpc.Dial(c.OtherServices.CartServiceUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
		return nil, err
	}
	return pb.NewCartServiceClient(cc), nil
}
func InitCourseServiceClient(c *config.Config) (course.CourseServiceClient, error) {
	// using WithInsecure() because no SSL running

	cc, err := grpc.Dial(c.OtherServices.CourseServiceUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
		return nil, err
	}
	return course.NewCourseServiceClient(cc), nil
}
