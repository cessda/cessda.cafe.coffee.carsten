pipeline{
    environment
	{
		project_name = "cessda-dev"
		module_name = "mgmt-coffeepot"
		image_tag = "eu.gcr.io/${project_name}/${module_name}:${env.BRANCH_NAME}-v${env.BUILD_NUMBER}"
	}

    agent any

    stages{
        stage("Build Docker Image"){
            steps{
                echo "Building Docker image using Dockerfile with tag"
                sh("docker build -t ${image_tag} .")
            }
        }
        stage('Push Docker image'){
            steps{
                echo 'Tag and push Docker image'
				sh("gcloud docker -- push ${image_tag}")
				sh("gcloud container images add-tag ${image_tag} eu.gcr.io/${project_name}/${module_name}:${env.BRANCH_NAME}-latest")
            }
        }
    }
}