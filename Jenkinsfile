pipeline{
    environment
	{
		project_name = "cessda-dev"
		module_name = "mgmt-coffeepot"
		image_tag = "eu.gcr.io/${project_name}/${module_name}:${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
        scannerHome = tool 'sonar-scanner'
	}

    agent any

    stages{
        stage('Run Test Suite'){
            agent{
                docker { image 'golang:latest' }
            }
            steps{
                echo "Running test suite"
                ln -s /go/src/carsten-coffee-api /root/workspace/${project_name}
                cd /go/src/carsten-coffee-api
                go get -u github.com/jstemmer/go-junit-report
                go get -u github.com/kardianos/govendor
                govendor sync
                go vet $(go list ./... | grep -v /vendor/)
                go test -v 2>&1 | go-junit-report > report.xml
            }
        }
        stage('Start Sonar scan') {
		    steps {
                withSonarQubeEnv('cessda-sonar') {
                    sh "${scannerHome}/bin/sonar-scanner"
                }
            }
        }
        stage("Get Sonar Quality Gate") {
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    waitForQualityGate abortPipeline: true
                }
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
            }
        }
        stage('Deploy Docker image'){
            steps{
                build job: '../cessda.coffeeapi.deployment/master', parameters: [string(name: 'DEPLOYMENT_VERSION', value: BUILD_NUMBER)], wait: false
            }
        }
    }
}
