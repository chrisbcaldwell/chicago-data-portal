# chicago-data-portal
Data lake serving selected data from the [City of Chicago Data Portal](https://data.cityofchicago.org/).  Course project for MSDS 432 Foundations of Data Engineering at Northwestern University.

## Backend Microservices
### Database Microservices
| Microservice | Platform      | Docker Image    |
| ------------ | --------------| ----------------|
| postgres     | PostgreSQL    | cbcaldwell/xxxx |
| chi-data     | Docker Volume | cbcaldwell/xxxx |
### Data Extraction Microservices
All data extraction microservices were built with Go.
| Microservice | Platform      | Docker Image    | Data Source URLs |

<table>
    <thead>
        <tr>
            <th>Microservice</th>
            <th>Docker Image</th>
            <th>Data Source URLs</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td rowspan=2>taxitrips</td>
            <td rowspan=2>cbcaldwell/xxxx</td>
            <td><a href="https://data.cityofchicago.org/Transportation/Taxi-Trips-2013-2023-/wrvz-psew/about_data">https://data.cityofchicago.org/Transportation/Taxi-Trips-2013-2023-/wrvz-psew/about_data</a></td>
        </tr>
        <tr>
            <td>L3 Name B</td>
        </tr>
        <tr>
            <td rowspan=2>L2 Name B</td>
            <td>L3 Name C</td>
        </tr>
        <tr>
            <td>L3 Name D</td>
        </tr>
    </tbody>
</table>

