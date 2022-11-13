// Carsten's CESSDA CAFE Coffee Machine
// Copyright CESSDA-ERIC 2019
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

pipeline{
    options {
        ansiColor('xterm')
        buildDiscarder logRotator(artifactNumToKeepStr: '5', numToKeepStr: '10')
    }

    environment
    {
        project_name = "${GCP_PROJECT}"
        product_name = "cafe"
        component_name = "coffee-carsten"
        image_tag = "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
        full_image_name = "${docker_repo}/${product_name}-${component_name}:${image_tag}"
        scannerHome = tool 'sonar-scanner'
    }

    agent any

    stages{
        stage('Run Test Suite'){
            agent{
                docker {
                    image 'golang:latest'
                    reuseNode true
                }
            }
            steps{
                echo "Running test suite"
                sh("ln -s $WORKSPACE /go/src/coffee-api")
                sh("cd /go/src/coffee-api && go get golang.org/x/lint/golint && go install golang.org/x/lint/golint")
                sh("cd /go/src/coffee-api && make test-ci")
            }
            post {
                always {
                    junit 'junit.xml'
                }
            }
        }
        stage('Run Sonar Scan') {
            steps {
                withSonarQubeEnv('cessda-sonar') {
                    sh "${scannerHome}/bin/sonar-scanner"
                }
            }
        }
        stage("Get Sonar Quality Gate") {
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    waitForQualityGate abortPipeline: false
                }
            }
        }
        stage("Build Docker Image"){
            steps{
                echo "Building Docker image using Dockerfile with tag ${full_image_name}"
                sh("docker build -t ${full_image_name} .")
            }
        }
        stage('Push Docker image'){
            steps{
                echo 'Tag and push Docker image'
                sh("gcloud auth configure-docker")
                sh("docker push ${full_image_name}")
                sh("gcloud container images add-tag ${full_image_name} ${docker_repo}/${product_name}-${component_name}:${env.BRANCH_NAME}-latest")
            }
        }
        stage('Deploy Docker image'){
            steps{
                build job: '../cessda.cafe.deployment/main', parameters: [string(name: 'image_tag', value: "${image_tag}"), string(name: 'component', value: "${component_name}")], wait: false
            }
        }
    }
}
