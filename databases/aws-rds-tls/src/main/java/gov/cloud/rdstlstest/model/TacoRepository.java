package gov.cloud.rdstlstest.model;

import org.springframework.data.repository.CrudRepository;

import gov.cloud.rdstlstest.model.Taco;

// This will be AUTO IMPLEMENTED by Spring into a Bean called tacoRepository
// CRUD refers Create, Read, Update, Delete

public interface TacoRepository extends CrudRepository<Taco, Integer> {

}
