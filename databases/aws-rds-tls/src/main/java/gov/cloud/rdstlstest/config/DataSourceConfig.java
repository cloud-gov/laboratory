package gov.cloud.rdstlstest.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;

import javax.sql.DataSource;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.jdbc.DataSourceBuilder;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.*;

@Configuration
public class DataSourceConfig {
    @Primary
    @Bean
    @ConfigurationProperties(prefix = "datasource")
    public DataSource datasource() throws Exception {
        String vcapServices = System.getenv("VCAP_SERVICES");
        Object servicesObject = new JSONParser().parse(vcapServices);
        JSONObject services = (JSONObject) servicesObject;
        JSONArray rdsServices = (JSONArray) services.get("aws-rds");
        JSONObject rdsServiceConfig = (JSONObject) rdsServices.get(0);
        JSONObject creds = (JSONObject) rdsServiceConfig.get("credentials");
        String url = (String) creds.get("host");
        String database = (String) creds.get("db_name");
        url = "jdbc:mysql://" + url + "/" + database + "?sslMode=VERIFY_CA&useSSL=true";

        String username = (String) creds.get("username");
        String password = (String) creds.get("password");

        return DataSourceBuilder.create()
           // .driverClassName("com.mysql.cj.jdbc.Driver")
            .url(url)
            .username(username)
            .password(password)
            .build();

    }
}
