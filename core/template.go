package core
const IndexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Branches Health</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/angular.js/1.6.1/angular.min.js" integrity="sha256-3SrLjPeRPa1ofM280r+OMcUjJZKLWJHr6SRtRu3dRb0=" crossorigin="anonymous"></script>
</head>
<body>
    <div ng-app="branchesHealth" ng-controller="app as $ctrl">
        <div ng-repeat="branch in $ctrl.branches">
            {{branch.Name}} - {{branch.IsMerged}} - {{branch.IsOutdated}}
        </div>
    </div>
</body>

<script type="application/javascript">

    function controller ($http) {
        this.branches = [];

        $http.get('/api/branches').then(response => {
            this.branches = response.data.branches;
        });
    }
    angular
        .module("branchesHealth", [])
        .controller("app", controller);

</script>

</html>`
