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

    <div ng-app="branchesHealth" ng-controller="app as $ctrl" id="content">

        <div class="block">
            <h1>Statistics</h1>
            <table class="statistics">
                <tr>
                    <td>Keep: {{$ctrl.info.keep}}</td>
                    <td>To Delete: {{$ctrl.info.toDelete}}</td>
                    <td>Old: {{$ctrl.info.old}}</td>
                </tr>
            </table>
        </div>

        <div class="block">

            <h1>Leaderboard</h1>
            <table class="statistics">
                <tr>
                    <td ng-repeat="t in ['keep', 'old', 'toDelete']">
                        <div ng-repeat="person in $ctrl.leaderboard[t]">
                            {{person.name}} - {{person[t]}}
                        </div>
                    </td>
                </tr>
            </table>

        </div>

        <div class="block">
            <h1>Branches</h1>
            <div ng-repeat="(repositoryName, branches) in $ctrl.repositories">
                <h2>{{repositoryName}}</h2>
                <table>
                    <thead>
                        <tr>
                            <th>Branch name</th>
                            <th colspan="2">Last changed by</th>
                            <th>Is merged?</th>
                            <th>Is outdated?</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr ng-repeat="branch in branches">
                            <td>{{branch.Name}}</td>
                            <td>{{branch.Author}}</td>
                            <td>{{branch.LastUpdated}}</td>
                            <td>{{branch.IsMerged ? 'yes' : 'no'}}</td>
                            <td>{{branch.IsOutdated ? 'yes' : 'no'}}</td>
                        </tr>
                    </tbody>
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

<style>

    /* http://meyerweb.com/eric/tools/css/reset/
       v2.0 | 20110126
       License: none (public domain)
    */

    html, body, div, span, applet, object, iframe,
    h1, h2, h3, h4, h5, h6, p, blockquote, pre,
    a, abbr, acronym, address, big, cite, code,
    del, dfn, em, img, ins, kbd, q, s, samp,
    small, strike, strong, sub, sup, tt, var,
    b, u, i, center,
    dl, dt, dd, ol, ul, li,
    fieldset, form, label, legend,
    table, caption, tbody, tfoot, thead, tr, th, td,
    article, aside, canvas, details, embed,
    figure, figcaption, footer, header, hgroup,
    menu, nav, output, ruby, section, summary,
    time, mark, audio, video {
        margin: 0;
        padding: 0;
        border: 0;
        font-size: 100%;
        font: inherit;
        vertical-align: baseline;
    }
    /* HTML5 display-role reset for older browsers */
    article, aside, details, figcaption, figure,
    footer, header, hgroup, menu, nav, section {
        display: block;
    }
    body {
        line-height: 1;
    }
    ol, ul {
        list-style: none;
    }
    blockquote, q {
        quotes: none;
    }
    blockquote:before, blockquote:after,
    q:before, q:after {
        content: '';
        content: none;
    }
    table {
        border-collapse: collapse;
        border-spacing: 0;
    }
</style>

<style>
    body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
        line-height: 1.3;
    }

    h1 {
        font-size: 2em;
        font-weight: 600;
        padding-top: .25em;
        padding-bottom: .5em;
    }

    h2 {
        font-size: 1.5em;
        font-weight: 600;
        padding-top: .5em;
        padding-bottom: .3em;
    }

    table {
        width: 100%;
    }

        table tr th {
            text-align: left;
            padding-bottom: .5em;
            color: #a2a2a2;
            font-size: .8em;
        }

        table.statistics tr td {
            width: 33.3%;
        }


    #content {
        width: 100%;
        max-width: 960px;
        margin: 0 auto;
    }

    .block {
        margin-bottom: 2em;
    }

</style>

</html>`
