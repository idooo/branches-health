package core
const IndexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Branches Health</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/angular.js/1.6.1/angular.min.js" integrity="sha256-3SrLjPeRPa1ofM280r+OMcUjJZKLWJHr6SRtRu3dRb0=" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/moment.js/2.17.1/moment.min.js" integrity="sha256-Gn7MUQono8LUxTfRA0WZzJgTua52Udm1Ifrk5421zkA=" crossorigin="anonymous"></script>
</head>
<body>
    <div ng-app="branchesHealth" ng-controller="app as $ctrl">

        <div class="statistics">
            <h1>Statistics</h1>
            <div>Keep: {{$ctrl.info.keep}}</div>
            <div>To Delete: {{$ctrl.info.toDelete}}</div>
            <div>Old: {{$ctrl.info.old}}</div>
        </div>

        <div class="leaderboard">

            <h1>Leaderboard</h1>
            <h2>Keep branches</h2>
            <div ng-repeat="person in $ctrl.leaderboard.keep">
                {{person.name}} - {{person.keep}}
            </div>
            <h2>Old branches</h2>
            <div ng-repeat="person in $ctrl.leaderboard.old">
                {{person.name}} - {{person.old}}
            </div>
            <h2>To delete</h2>
            <div ng-repeat="person in $ctrl.leaderboard.toDelete">
                {{person.name}} - {{person.toDelete}}
            </div>

        </div>

        <div class="branches">
            <h1>Branches</h1>
            <div ng-repeat="(repositoryName, branches) in $ctrl.repositories">
                <h2>{{repositoryName}}</h2>
                <table>
                    <tr ng-repeat="branch in branches">
                        <td>{{branch.Name}}</td>
                        <td>{{branch.Author}}</td>
                        <td>{{branch.LastUpdated}}</td>
                        <td>{{branch.IsMerged}}</td>
                        <td>{{branch.IsOutdated}}</td>
                    </tr>
                </table>
            </div>
        </div>
    </div>
</body>

<script type="application/javascript">

    function controller ($http) {
        this.repositories = [];
        this.info = {
            keep: 0,
            toDelete: 0,
            old: 0
        };
        this.people = {};
        this.leaderboard = {
            keep: [],
            toDelete: [],
            old: []
        };

        $http.get('/api/branches').then(response => {
            const branches = response.data.branches;

            branches.forEach(branch => {
                if (this.people[branch.Author] === undefined) {
                    this.people[branch.Author] = {
                        keep: 0,
                        toDelete: 0,
                        old: 0
                    };
                }
                if (branch.IsMerged) {
                    this.people[branch.Author].toDelete++;
                    this.info.toDelete++;
                }
                else if (branch.IsOutdated) {
                    this.people[branch.Author].old++;
                    this.info.old++;
                }
                else {
                    this.people[branch.Author].keep++;
                    this.info.keep++;
                }
            });

            {
                const people = [];
                Object.keys(this.people)
                    .forEach(name => people.push(Object.assign({name}, this.people[name])));

                this.leaderboard.toDelete = people.sort((a, b) => b.toDelete - a.toDelete);
                this.leaderboard.old = people.sort((a, b) => b.old - a.old);
                this.leaderboard.keep = people.sort((a, b) => b.keep - a.keep);
            }

        });

        $http.get('/api/repositories').then(response => {
            this.repositories = response.data;
            Object.keys(this.repositories).forEach(repositoryName => {
                this.repositories[repositoryName].forEach(branch => {
                    branch.LastUpdated = moment(branch.LastUpdated).fromNow()
                })
            })
        })
    }

    angular
        .module("branchesHealth", [])
        .controller("app", controller);

</script>

</html>`
