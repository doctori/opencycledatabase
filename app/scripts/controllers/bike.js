'use strict';

/**
 * @ngdoc function
 * @name openBicycleDatabaseApp.controller:BikecontrollerCtrl
 * @description
 * # BikecontrollerCtrl
 * Controller of the openBicycleDatabaseApp
 */

 function bikeListCtrl(Bike, $scope, $mdSidenav, $mdBottomSheet) {
    var self = this;
    self.selected     = null;
    self.bikes        = [ ];
    $scope.itemsPerPage = 30;
    $scope.searchString = "";
    $scope.currentPage = 0;
    $scope.total = 9999;
    Bike.query(function (bikes){
        self.bikes = [].concat(bikes)
        self.selected = bikes[0];
      });
      
    self.selectBike = selectBike;
    self.toggleList = toggleUsersList;

    self.prevPage = function(){
      if($scope.currentPage > 0){
        $scope.currentPage--;
      }
    }
    self.nextPage = function(){
      if($scope.currentPage < self.pageCount() - 1) {
        $scope.currentPage++;
      }
    }
    self.prevPageDisabled = function() {
      return $scope.currentPage === 0 ? "disabled" : "";
    };
    self.nextPageDisabled = function() {
      return $scope.currentPage === self.pageCount() - 1 ? "disabled" : "";
    };
    self.pageCount = function() {
      return Math.ceil($scope.total/$scope.itemsPerPage);
    };
    $scope.$watch("currentPage", function(newValue,oldValue){
      Bike.getRange({page:newValue,per_page:$scope.itemsPerPage,search_string:$scope.searchString},function(bikes){
        self.bikes = [].concat(bikes)
      });
    })
    $scope.$watch("searchString",function(newValue,oldValue){
      if(newValue != ""){
        Bike.getRange({page:$scope.currentPage,per_page:$scope.itemsPerPage,search_string:newValue},function(bikes){
          self.bikes = [].concat(bikes)
        });
      }
    })
    self.remove = function(index,bike){
      console.log(index);
      console.log(bike);
      bike.$remove(function(){
        self.bikes = Bike.query();
      });
    };



   /**
   * Select the current bike
   * @param menuId
   */
    function selectBike ( bike ) {
      console.log(" BIKE SELECTED");
      console.log(bike)
      bike = angular.isNumber(bike) ? self.bikes[bike] : bike;
      Bike.get({id:bike.ID},function(bike){
        self.selected = bike;
      });
    }
   /**
   * Hide or Show the 'left' sideNav area
   */
    function toggleUsersList() {
      $mdSidenav('left').toggle();
    }

  }
 

angular.module('openBicycleDatabaseApp')
  .controller('BikeListCtrl',['Bike','$scope','$mdSidenav','$mdBottomSheet',bikeListCtrl ])
 /* .controller('BikeDetailCtrl', function ($scope, $routeParams,Bike) {
    Bike.get({id:$routeParams.bikeID},function(bike){
      console.log(bike);
      $scope.bike = bike;
    });
    $scope.bikeSave = function(){
      console.log($scope.bike);
      return $scope.bike.$save();
    };
    $scope.newComponentForm = function(){
      console.log("instantiating new Component to insert into the Bike");
      $scope.newComponent = {};
    };
    $scope.getImgURL = function(bike){
      if(bike.Image){
        var bikeURL= 'http://192.168.198.130:8080/image/'+bike.Image  
      }else{
        var bikeURL= "http://monculsurlacomode.com/"
      }
      
      console.log(bikeURL)
      return bikeURL
    }
    $scope.saveComponent = function(component){
      console.log("adding :");
      console.log(component);
      if($scope.bike.Components == null){
        $scope.bike.Components = [];
      }
      $scope.bike.Components.push(component);
      $('#new-component').modal('toggle');
      return $scope.bike.$save();
    };
    $scope.nameLabel = "Name";
    $scope.brandLabel = "Brand";
    $scope.yearLabel = "Year";
    $scope.descriptionLabel = "Description";
    $scope.componentNameLabel = "Component Name";
  })
  .controller('NewBikeCtrl',['$scope','Upload','$timeout','Bike', function($scope,Upload,$timeout, Bike){
    //Label Def
    $scope.nameLabel  = 'Name';
    $scope.brandLabel = 'Brand';
    $scope.yearLabel  = 'Year';
    $scope.namePlaceHolder = 'Bike Name';
    // Errors Def
    $scope.errorBrandRequired = 'The Brand Name is Required';
    $scope.errorYearRequired = 'The Year of the model is Required';
    $scope.errorNameRequired = 'The Name Of the Model is Required';
    $scope.errorYearPattern = ' The Year must be a Digit';
    // Misc
    $scope.validate = 'Submit';
    $scope.$watch('files', function () {
        $scope.upload($scope.files);
    });
    $scope.$watch('file', function () {
        if ($scope.file != null) {
            $scope.files = [$scope.file]; 
        }
    });
    $scope.log = '';
    $scope.imageId = 0;
    $scope.upload = function (files) {
          if (files && files.length) {
              for (var i = 0; i < files.length; i++) {
                var file = files[i];
                if (!file.$error) {
                  Upload.upload({
                      url: 'http://192.168.198.130:8080/image',
                      data: {
                        username: $scope.username,
                        file: file  
                      }
                  }).then(function (resp) {
                      $timeout(function() {
                          $scope.log = 'file: ' +
                          resp.config.data.file.name +
                          ', Response: ' + JSON.stringify(resp.data) +
                          '\n' + $scope.log;
                          $scope.imageId = resp.data.ID
                      });
                  }, null, function (evt) {
                      var progressPercentage = parseInt(100.0 *
                          evt.loaded / evt.total);
                      $scope.log = 'progress: ' + progressPercentage + 
                        '% ' + evt.config.data.file.name + '\n' + 
                        $scope.log;
                  });
                }
              }
          }
      };    
    $scope.save = function(bike){
      console.log(bike);
      if ($scope.imageId != 0){
        bike.image = $scope.imageId;  
      }
      Bike.save(bike);
    };
  }]);
  */
