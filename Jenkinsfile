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
        buildDiscarder logRotator(artifactNumToKeepStr: '5', numToKeepStr: '10')
    }

    environment
    {
        project_name = "${GCP_PROJECT}"
        product_name = "cafe"
        module_name = "coffeepot"
        image_tag = "${docker_repo}/${product_name}-${module_name}:${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
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
                sh("ln -s $WORKSPACE /go/src/carsten-coffee-api")
                sh("cd /go/src/carsten-coffee-api && ./run-tests.sh")
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
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        stage("Build Docker Image"){
            steps{
                echo "Building Docker image using Dockerfile with tag ${image_tag}"
                sh("docker build -t ${image_tag} .")
            }
        }
        stage('Push Docker image'){
            steps{
                echo 'Tag and push Docker image'
                sh("gcloud auth configure-docker")
                sh("docker push ${image_tag}")
                sh("gcloud container images add-tag ${image_tag} ${docker_repo}/${product_name}-${module_name}:${env.BRANCH_NAME}-latest")
            }
        }
        stage('Deploy Docker image'){
            steps{
                build job: '../cessda.cafe.deployment/master', parameters: [string(name: 'coffeepot_image_tag', value: "${image_tag}"), string(name: 'module', value: 'coffeepot')], wait: false
            }
        }
    }
}
