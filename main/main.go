package main

import (
	"github.com/jbowles/wordlab"
)

func main() {
	var root_errorfp = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data/"
	var root_datafp = "datasets/"
	wordlab.BuildHotelProviderDataKnn(root_errorfp, root_datafp)
	//wordlab.BuildHotelProviderDataKnnAmit(root_errorfp, root_datafp)
}

/*
import (
	"fmt"
	//"github.com/jbowles/wordlab"
	//"strings"
	"github.com/sjwhitworth/golearn/base"
)

func main() {

			newInst := base.NewDenseInstances()
			fmt.Println("building new attributes... ")
			attrs := make([]base.Attribute, 2) //number of attributes
			attrs[0] = base.NewFloatAttribute("FitValue")
			attrs[1] = new(base.CategoricalAttribute)
			attrs[1].SetName("Category")

			fmt.Printf("defined attributes and categorical... %v\n\n", newInst)

			newSpecs := make([]base.AttributeSpec, len(attrs))
			for i, a := range attrs {
				fmt.Printf("i: %v\t a: %v\n\n\n", i, a)
				newSpecs[i] = newInst.AddAttribute(a)
			}
			fmt.Printf("added attributes... %v\n\n", newInst)

		neighbors := 12
		distance := "manhattan"
		training := "datasets/wordlab_knn_hotel_error_test.csv"
		knnClassifier, testData := wordlab.InitKnnClassifier(neighbors, distance, training)

		//dense_copy := base.NewDenseCopy(testData)
		//fmt.Printf("copy of dense instance knnClasifier: %v\n", dense_copy)
		//sent := "We regret to inform you your credit card was declined."

		fmt.Printf("%v", knnClassifier)
		fmt.Printf("%v \n\n\n", testData)
			//fmt.Printf("*** Prections:: %v\n\n", knnClassifier.Predict(testData))

				var root_errorfp = "/Users/jbowles/x/training_data/partner_fusion_trained_errors/training_data/"
				var root_datafp = "data/"
				wordlab.BuildHotelProviderDataKnn(root_errorfp, root_datafp)
}

	FloatAttribute(FitValue)
*	CategoricalAttribute("Category", [availability_error booking_error cancel_error cancel_forbidden_error credit_data_error credit_decline_error credit_service_error unexpected_response_error])
*/
