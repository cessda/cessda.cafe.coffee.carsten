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
            agent{
                docker { image 'golang:latest' }
            }
            steps{
                echo "Running test suite"
                sh("run-tests.sh")
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
    post {
        always {
            junit 'report.xml'
        }
    }
}