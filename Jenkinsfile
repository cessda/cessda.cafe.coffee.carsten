pipeline{
    environment
	{
		project_name = "cessda-dev"
		module_name = "mgmt-coffeepot"
		image_tag = "eu.gcr.io/${project_name}/${module_name}:${env.BRANCH_NAME}-v${env.BUILD_NUMBER}"
	}

    agent any

    stages{
        stage('Run Test Suite'){
            steps{
                echo "Running test suite"
                sh("mkdir -p /go/src/nsd-utvikling /go/src/_/builds")
                sh("cp -r . /go/src/nsd-utvikling/carsten-coffee-api")
                dir('/go/src/nsd-utvikling/carsten-coffee-api')
                sh("go get -u github.com/jstemmer/go-junit-report")
                sh("go get -u github.com/kardianos/govendor")
                sh("govendor sync")
                sh("go vet $(go list ./... | grep -v /vendor/)")
                sh("go test -v 2>&1 | go-junit-report > report.xml")
            }
        }
        stage("Build Docker Image"){
            steps{
                echo "Building Docker image using Dockerfile with tag"
                sh("gcloud builds submit --tag ${image_tag} .")
            }
        }
        stage('Push Docker image'){
            steps{
                echo 'Tag and push Docker image'
				sh("gcloud container images add-tag ${image_tag} eu.gcr.io/${project_name}/${module_name}:${env.BRANCH_NAME}-latest")
                sh("gcloud container images add-tag ${image_tag} eu.gcr.io/${project_name}/${module_name}:latest")
            }
        }
    }
}