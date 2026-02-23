pipeline {
  agent any
  
  options {
    disableConcurrentBuilds()
  }
  
  stages {
    stage('install dependency') {
      steps {
        sh 'make install_dependency_frontend'
      }
    }

    stage('code analysis') {
      parallel {
        stage('code analysis frontend') {
          steps {
            sh 'make code_analysis_frontend'
          }
        }

        stage('code analysis backend') {
          steps {
            sh 'make code_analysis_backend'
          }
        }

      }
    }

    stage('unit test backend') {
      steps {
        sh 'make backend_unit_test'
        junit 'store-service/*.xml'
      }
    }

    stage('setup test fixtures') {
      steps {
        sh 'make setup_test_fixtures'
        sh 'sleep 8'
      }
    }

    stage('run integration test') {
      steps {
        sh 'make backend_integration_test'
      }
    }

    stage('build') {
      parallel {
        stage('build frontend') {
          steps {
            sh 'make build_frontend'
          }
        }

        stage('build backend') {
          steps {
            sh 'make build_backend'
          }
        }

      }
    }

    stage('run ATDD') {
      steps {
        sh 'make start_test_suite_grid'
        // sh 'make run_newman'
        sh 'make run_robot_grid'
        // robot outputPath: './atdd/ui', passThreshold: 100.0
        junit 'atdd/ui/reports/*.xml'
        sh 'make stop_test_suite'
      }
    }

    stage('trigger deployment') {
      steps {
        script {
          // Trigger deployment pipeline with current BUILD_NUMBER as IMAGE_TAG
          build job: 'sck-online-store-deploy', 
            parameters: [
              string(name: 'IMAGE_TAG', value: "${BUILD_NUMBER}"),
              string(name: 'ENVIRONMENT', value: 'development'),
              booleanParam(name: 'DEPLOY_ALL', value: true)
            ],
            wait: false
        }
      }
    }
    
  }
  
  post {
    always {
      sh 'make stop_test_suite'
      sh 'docker volume prune -f'
    }

  }
}
