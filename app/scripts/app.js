'use strict';

/**
 * @ngdoc overview
 * @name openBicycleDatabaseApp
 * @description
 * # openBicycleDatabaseApp
 *
 * Main module of the application.
 */
angular
  .module('openBicycleDatabaseApp', [
    'bikeService',
    'componentService',
    'ngCookies',
    'ngMaterial',
    'ngResource',
    'ngRoute',
    'gettext',
    'xeditable',
    'ngFileUpload'
  ])
  .run(function(gettextCatalog){
    gettextCatalog.setCurrentLanguage('fr_FR');
    gettextCatalog.debug = true;
  })
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/bike.html',
        controller: 'BikeListCtrl',
        controllerAs: 'bikeCtrl'
      })
      .when('/bike/:bikeID',{
        templateUrl: 'views/bike-detail.html',
        controller: 'BikeDetailCtrl',
        controllerAs: 'bike'
      })
      .when('/new-bike',{
        templateUrl: 'views/new-bike.html',
        controller: 'NewBikeCtrl',
        controllerAs: 'bike'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl',
        controllerAs: 'about'
      })
      .when('/components', {
        templateUrl: 'views/components.html',
        controller: 'ComponentListCtrl',
        controllerAs: 'component'
      })

      .otherwise({
        redirectTo: '/'
      });
  });