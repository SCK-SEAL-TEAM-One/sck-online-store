pipeline {
  agent any
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
        sh 'make start_test_suite'
        // sh 'make run_newman'
        sh 'make run_robot'
        robot outputPath: './atdd/ui', passThreshold: 100.0
        sh 'make stop_test_suite'
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
