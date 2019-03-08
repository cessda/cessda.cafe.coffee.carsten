pipeline{
    environment
	{
		project_name = "cessda-dev"
		module_name = "mgmt-coffeepot"
		image_tag = "eu.gcr.io/${project_name}/${module_name}:${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
	}

    agent any

    stages{
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
                build job: '../cessda.coffeeapi.deployment/master', parameters: [string(name: 'DEPLOYMENT_VERSION', value: "{env.BUILD_NUMBER}")], wait: false
            }
        }
    }
}