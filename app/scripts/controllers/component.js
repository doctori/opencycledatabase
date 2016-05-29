'use strict';

/**
 * @ngdoc function
 * @name openBicycleDatabaseApp.controller:ComponentcontrollerCtrl
 * @description
 * # ComponentcontrollerCtrl
 * Controller of the openBicycleDatabaseApp
 */
function componentListCtrl ($scope,Component){
	var self = this;
 	self.selected     = null;
    self.components        = [ ];
     $scope.itemsPerPage = 30;
    $scope.searchString = "";
    $scope.currentPage = 0;
    $scope.total = 9999;
    Component.query(function (components){
        self.components = [].concat(components)
        self.selected = components[0];
        console.log(self.components)
      });
     self.selectComponent = selectComponent;

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
      Component.getRange({page:newValue,per_page:$scope.itemsPerPage,search_string:$scope.searchString},function(components){
        self.components = [].concat(components)
      });
    })
    $scope.$watch("searchString",function(newValue,oldValue){
      if(newValue != ""){
        Component.getRange({page:$scope.currentPage,per_page:$scope.itemsPerPage,search_string:newValue},function(components){
          self.components = [].concat(components)
        });
      }
    })
    /*self.remove = function(index,component){
      console.log(index);
      console.log(component);
      component.$remove(function(){
        self.components = Component.query();
      });
    };*/
      /**
   * Select the current component
   * @param menuId
   */
    function selectComponent ( component ) {
      console.log(" COMPONENT SELECTED");
      console.log(component)
      component = angular.isNumber(component) ? self.components[component] : component;
      Component.get({id:component.ID},function(component){
        self.selected = component;
      });
    }
}
angular.module('openBicycleDatabaseApp')
  .controller('ComponentListCtrl',['$scope','Component',componentListCtrl ])
  
 