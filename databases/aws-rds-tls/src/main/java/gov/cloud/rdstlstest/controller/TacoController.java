package gov.cloud.rdstlstest.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import gov.cloud.rdstlstest.model.*;

@Controller
@RequestMapping(path="/taco")
public class TacoController {
	@Autowired
	private TacoRepository tacoRepository;

	@PostMapping(path="/add")
	public @ResponseBody String addNewTaco (@RequestParam String topping) {
		Taco n = new Taco();
		n.setTopping(topping);
		tacoRepository.save(n);
		return "Saved";
	}

	@GetMapping(path="/all")
	public @ResponseBody Iterable<Taco> getAllTacos() {
		return tacoRepository.findAll();
	}
}
