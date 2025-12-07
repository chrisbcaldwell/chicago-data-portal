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
            <td><a href="https://data.cityofchicago.org/Transportation/Taxi-Trips-2024-/ajtu-isnz/about_data">https://data.cityofchicago.org/Transportation/Taxi-Trips-2024-/ajtu-isnz/about_data</a></td>
        </tr>
        <tr>
            <td rowspan=3>rideshare</td>
            <td rowspan=3>cbcaldwell/xxxx</td>
            <td><a href="https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips-2018-2022-/m6dm-c72p/about_data">https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips-2018-2022-/m6dm-c72p/about_data</a></td>
        </tr>
        <tr>
            <td><a href="https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips-2023-2024-/n26f-ihde/about_data">https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips-2023-2024-/n26f-ihde/about_data</a></td>
        </tr>
        <tr>
            <td><a href="https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips-2025-/6dvr-xwnh/about_data">https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips-2025-/6dvr-xwnh/about_data</td>
        </tr>
    </tbody>
</table>

